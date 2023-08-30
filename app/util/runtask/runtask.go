package runtask

import (
	"fmt"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"pear-admin-go/app/core/log"
	"pear-admin-go/app/core/scli"
	"pear-admin-go/app/dao"
	"pear-admin-go/app/util/file"
	"pear-admin-go/app/util/sysos"

	"pear-admin-go/app/global/e"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/check"
	"pear-admin-go/app/util/pool"
	"strings"
	"sync/atomic"
	"time"
)

type RunTask struct {
	task         model.Task
	sourceClient *sftp.Client // 源服务器连接
	dstClient    *sftp.Client // 目标服务器连接
	fp           *pool.Pool   // chan pool
	counter      uint64       // 文件传输计数器
}

func NewRunTask(task model.Task) *RunTask {
	return &RunTask{task: task, fp: pool.NewPool(e.MaxPool), counter: 0}
}

func (this *RunTask) SetSourceClient() *RunTask {
	if this.task.SourceType == e.Local {
		return this
	}
	c, err := this.getClient(this.task.SourceServer)
	if err != nil {
		log.Instance().Error("SetSourceClient.getClient", zap.Error(err))
		return this
	}
	this.sourceClient = c
	return this
}

func (this *RunTask) SetDstClient() *RunTask {
	if this.task.DstType == e.Local {
		return this
	}
	c, err := this.getClient(this.task.DstServer)
	if err != nil {
		log.Instance().Error("SetSourceClient.getClient", zap.Error(err))
		return this
	}
	this.dstClient = c
	return this
}

func (this *RunTask) getClient(sid int) (*sftp.Client, error) {
	server, err := dao.NewTaskServerDaoImpl().FindOne(sid)
	if err != nil {
		return nil, err
	}
	c, err := scli.Instance(*server)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (this *RunTask) Run() {
	if this.task.SourceType == e.Remote && this.task.DstType == e.Remote {
		this.RunR2R()
	} else if this.task.SourceType == e.Remote && this.task.DstType == e.Local {
		this.RunR2L()
	} else if this.task.SourceType == e.Local && this.task.DstType == e.Remote {
		this.RunL2R()
	}
	err := dao.NewTaskDaoImpl().Update(this.task, map[string]interface{}{"task_file_num": this.counter})
	if err != nil {
		log.Instance().Error("Run.Update", zap.Error(err))
	}
}

func (this *RunTask) RunR2R() {
	this.WalkRemotePath(this.task.SourcePath, e.ToRemote)
}

func (this *RunTask) RunR2L() {
	this.WalkRemotePath(this.task.SourcePath, e.ToLocal)
}

func (this *RunTask) WalkRemotePath(dirPath string, runType int) {
	globPath := pathJoin(dirPath)
	files, err := this.sourceClient.Glob(globPath)
	if err != nil {
		log.Instance().Error("WalkRemotePath.this.sourceClient.Glob", zap.Error(err))
		return
	}
	for _, v := range files {
		stat, err := this.sourceClient.Stat(v)
		if err != nil {
			log.Instance().Error("WalkRemotePath.this.sourceClient.Stat", zap.Error(err))
			continue
		}
		if stat.IsDir() {
			if runType == e.ToRemote {
				dname := string([]rune(strings.ReplaceAll(v, this.task.SourcePath, ""))[1:])
				err = this.MkRemotedir(dname)
				if err != nil {
					log.Instance().Error("WalkRemotePath.MkRemotedir", zap.Error(err))
					return
				}
			}
			this.WalkRemotePath(v, runType)
		} else {
			this.fp.Add(1)
			atomic.AddUint64(&this.counter, 1)
			go func(v string, size int64) {
				defer this.fp.Done()
				if runType == e.ToLocal {
					err = this.RemoteSendLocal(v, size)
				} else if runType == e.ToRemote {
					err = this.RemoteSendRemote(v, size)
				}
				if err != nil {
					log.Instance().Error("WalkRemotePath.RemoteToLocal", zap.Error(err))
				}
			}(v, stat.Size())
		}
	}
}

func (this *RunTask) RemoteSendRemote(fname string, fsize int64) error {
	newName := strings.ReplaceAll(fname, this.task.SourcePath, "")
	rf := path.Join(this.task.DstPath, newName) // 文件在目标服务器的路径及名称
	has, err := this.dstClient.Stat(rf)
	if err == nil && (has.Size() == fsize) {
		log.Instance().Debug(fmt.Sprintf("文件%s已存在", rf))
		return nil
	}
	err = this.dstClient.MkdirAll(this.task.DstPath)
	if err != nil {
		return err
	}
	err = this.dstClient.Chmod(this.task.DstPath, os.ModePerm)
	if err != nil {
		return err
	}
	srcFile, err := this.sourceClient.Open(fname)
	if err != nil {
		log.Instance().Error("RemoteToLocal.sourceClient.Open", zap.Error(err))
		return err
	}
	defer srcFile.Close()

	dstFile, err := this.dstClient.Create(rf) // 如果文件存在，create会清空原文件 openfile会追加
	if err != nil {
		log.Instance().Error("RemoteSendRemote.this.dstClient.Create", zap.Error(err))
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 10000)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[:n]) // 读多少 写多少
	}
	log.Instance().Info(fmt.Sprintf("【%s】传输完毕", fname))
	err = dao.NewTaskLogDaoImpl().Insert(model.TaskLog{
		TaskId:     this.task.Id,
		ServerId:   this.task.DstServer,
		SourcePath: fname,
		DstPath:    rf,
		Size:       fsize,
		CreateTime: time.Now(),
	})
	if err != nil {
		log.Instance().Error("RemoteSendRemote.dao.NewTaskLogDaoImpl.Insert", zap.Error(err))
		return err
	}
	return nil
}

// 远端->本地 使用 sourceClient
func (this *RunTask) RemoteSendLocal(fname string, fsize int64) error { // 本地文件夹
	dstFile := path.Join(this.task.DstPath, strings.ReplaceAll(fname, this.task.SourcePath, "")) // 需要保存的本地文件地址
	has, err := check.CheckFile(dstFile)                                                         // 是否已存在
	if err != nil {
		log.Instance().Error("RemoteToLocal.CheckFile", zap.Error(err))
		return err
	}
	if has != nil && has.Size() == fsize {
		log.Instance().Debug(fmt.Sprintf("文件%s已存在", dstFile))
		return nil
	}
	dir, _ := path.Split(dstFile)
	err = file.IsNotExistMkDir(dir)
	if err != nil {
		log.Instance().Error("RemoteToLocal.IsNotExistMkDir", zap.Error(err))
		return err
	}

	srcFile, err := this.sourceClient.Open(fname)
	if err != nil {
		log.Instance().Error("RemoteToLocal.sourceClient.Open", zap.Error(err))
		return err
	}
	defer srcFile.Close()
	lf, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer lf.Close()

	if _, err = srcFile.WriteTo(lf); err != nil {
		return err
	}
	log.Instance().Info(fmt.Sprintf("【%s】传输完毕", srcFile.Name()))

	err = dao.NewTaskLogDaoImpl().Insert(model.TaskLog{
		TaskId:     this.task.Id,
		ServerId:   this.task.DstServer,
		SourcePath: fname,
		DstPath:    dstFile,
		Size:       fsize,
		CreateTime: time.Now(),
	})
	if err != nil {
		log.Instance().Error("RemoteSendRemote.dao.NewTaskLogDaoImpl.Insert", zap.Error(err))
		return err
	}
	return nil
}

func pathJoin(p string) (np string) {
	if strings.HasSuffix(p, "/") == false {
		p = p + "/"
	}
	if sysos.IsWindows() {
		np = strings.ReplaceAll(path.Join(p, "*"), "\\", "/")
	} else {
		np = path.Join(p, "*")
	}
	return
}

// 本地->远端
func (this *RunTask) RunL2R() {
	_ = filepath.Walk(this.task.SourcePath, func(v string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Instance().Error("RunL2R.Walk.err", zap.Error(err))
			return nil
		}
		if info == nil {
			log.Instance().Error("RunL2R.Walk.info Is Nil")
			return nil
		}
		stat, err := os.Stat(v)
		if err != nil {
			log.Instance().Error("RunL2R.os.Stat", zap.Error(err))
			return nil
		}
		if stat.IsDir() {
			if v == this.task.SourcePath {
				return nil
			}
			dname := string([]rune(strings.ReplaceAll(v, this.task.SourcePath, ""))[1:])
			err = this.MkRemotedir(dname)
			if err != nil {
				log.Instance().Error("RunL2R.rm.Mkdir", zap.Error(err))
				return nil
			}
		} else {
			this.fp.Add(1)
			atomic.AddUint64(&this.counter, 1)
			go func(v string, size int64) {
				defer this.fp.Done()
				err = this.LocalSend(v, size)
				if err != nil {
					log.Instance().Error("WalkPath.LocalToRemote", zap.Error(err))
				}
			}(v, stat.Size())
		}
		return nil
	})
}

func (this *RunTask) LocalSend(fname string, fsize int64) error {
	if sysos.IsWindows() {
		this.task.SourcePath = strings.ReplaceAll(this.task.SourcePath, "\\", "/")
		this.task.DstPath = strings.ReplaceAll(this.task.DstPath, "\\", "/")
		fname = strings.ReplaceAll(fname, "\\", "/")
		fname = strings.ReplaceAll(fname, this.task.SourcePath, "")
	}
	rf := path.Join(this.task.DstPath, fname) // 文件在服务器的路径及名称
	has, err := this.dstClient.Stat(rf)
	if err == nil && (has.Size() == fsize) {
		log.Instance().Debug(fmt.Sprintf("文件%s已存在", rf))
		return nil
	}
	err = this.dstClient.MkdirAll(this.task.DstPath)
	if err != nil {
		return err
	}
	err = this.dstClient.Chmod(this.task.DstPath, os.ModePerm)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(path.Join(this.task.SourcePath, fname))
	if err != nil {
		log.Instance().Error("源文件无法读取", zap.Error(err))
		return err
	}
	defer srcFile.Close()
	dstFile, err := this.dstClient.Create(rf) // 如果文件存在，create会清空原文件 openfile会追加
	if err != nil {
		log.Instance().Error("this.dstClient.Create", zap.Error(err))
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 10000)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[:n]) // 读多少 写多少
	}
	err = dao.NewTaskLogDaoImpl().Insert(model.TaskLog{
		TaskId:     this.task.Id,
		ServerId:   this.task.DstServer,
		SourcePath: path.Join(this.task.SourcePath, fname),
		DstPath:    rf,
		Size:       fsize,
		CreateTime: time.Now(),
	})
	if err != nil {
		return err
	}
	log.Instance().Info(fmt.Sprintf("【%s】传输完毕", fname))
	return nil
}

func (this *RunTask) MkRemotedir(p string) error {
	p = check.CheckWinPath(p)
	dst := path.Join(this.task.DstPath, p)
	err := this.dstClient.MkdirAll(dst)
	if err != nil {
		log.Instance().Error("MkRemotedir", zap.Error(err))
		return err
	}
	return nil
}
