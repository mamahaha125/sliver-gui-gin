// @Title       filelist.go
// @Description //TODO
// @Author      mamahaha125
// @Data        2023/8/24 11:23
package service

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	"path/filepath"
	"pear-admin-go/app/core/cmd"
	rpcs "pear-admin-go/app/core/rpc"
	"strings"
)

type FileList struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Mode    string `json:"mode"`
	ModTime int64  `json:"modTime"`
	Size    int64  `json:"size"`
	Path    string `json:"path"`
	Id      string `json:"id"`
}

// @Title GetFiles
// @Description 获取文件列表
// @Author mamahaha125 2023-08-28 21:38:27
// @Update mamahaha125 2023-08-28 21:38:27
// @Param path
// @Return []FileList
// @Return string
func GetFiles(path string) ([]FileList, string) {

	//clientName := "new@192.168.200.101 (6b48954e9c0de3e1)"

	//rpc, con, _ := ConnectRpc(clientName)
	//SetCon(con, rpc)

	if ok := strings.Contains(path, ":"); !ok {
		path = strings.Replace(path, "\\", "/", -1)
	}

	con := rpcs.GetInstance().GetCon()
	cmds := cmd.GetInstance().GetCmd()

	ls, _ := con.Rpc.Ls(context.Background(), &sliverpb.LsReq{
		Request: con.ActiveTarget.Request(cmds),
		Path:    path,
	})
	var List []FileList
	var File FileList
	var num int
	for _, v := range ls.Files {
		if v.IsDir {
			File.Type = "dir"
		}
		num = num + 1
		File.Id = fmt.Sprintf("%03d", num)
		File.Title = v.Name
		File.Size = v.Size
		File.Mode = v.Mode
		File.ModTime = v.ModTime
		File.Path = filepath.Join(ls.Path, v.Name)
		File.Type = "notdir"
		List = append(List, File)
	}
	if path == "." || path == "" {
		pwd, err := con.Rpc.Pwd(context.Background(), &sliverpb.PwdReq{
			Request: con.ActiveTarget.Request(cmds),
		})
		if err != nil {
			return nil, ""
		}
		return List, pwd.Path
	}
	return List, path
}
