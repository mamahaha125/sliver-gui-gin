package service

import (
	"encoding/json"
	"pear-admin-go/app/core/cache"
	dao2 "pear-admin-go/app/dao"
	e2 "pear-admin-go/app/global/e"
	"pear-admin-go/app/global/request"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
	"strings"
	"time"
)

func SiteEditService(f request.SiteConfForm) error {
	if strings.HasPrefix(f.WebUrl, "http://") == false { // http前缀校验
		f.WebUrl = "http://" + f.WebUrl
	}
	if strings.HasSuffix(f.WebUrl, "/") == false { // 结尾校验
		f.WebUrl = f.WebUrl + "/"
	}
	if f.WebUrl != "" && strings.HasPrefix(f.LogoUrl, f.WebUrl) == false { // logo地址校验
		f.LogoUrl = f.WebUrl + f.LogoUrl
	}
	str, err := json.Marshal(f)
	if err != nil {
		return err
	}
	if f.ID == 0 { // 新增
		err = dao2.NewSysConfDaoImpl().Insert(model.SysConf{
			Type:      model.SysSiteConf,
			Info:      gconv.String(str),
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	} else {
		err = dao2.NewSysConfDaoImpl().Update(f.ID, model.SysConf{
			Type:      model.SysSiteConf,
			Info:      gconv.String(str),
			Status:    1,
			UpdatedAt: time.Now(),
		})
	}
	if err != nil {
		return err
	}
	return nil
}

func MailEditService(f request.MailConfForm) error {
	str, err := json.Marshal(f)
	if err != nil {
		return err
	}
	if f.ID == 0 { // 新增
		err = dao2.NewSysConfDaoImpl().Insert(model.SysConf{
			Type:      model.SysMailConf,
			Info:      gconv.String(str),
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	} else {
		err = dao2.NewSysConfDaoImpl().Update(f.ID, model.SysConf{
			Type:      model.SysMailConf,
			Info:      gconv.String(str),
			Status:    1,
			UpdatedAt: time.Now(),
		})
	}
	if err != nil {
		return err
	}
	return nil
}

func GetMailTestConf() (testMail model.MailTest) {
	mail, has := cache.Instance().Get(e2.TestMailConf)
	if has {
		testMail = mail.(model.MailTest)
	}
	return testMail
}

func GetSiteConf() (site model.SiteConf, sysID uint) {
	sysInfo, _ := dao2.NewSysConfDaoImpl().FindBySysType(model.SysSiteConf)
	if sysInfo.ID > 0 {
		json.Unmarshal([]byte(sysInfo.Info), &site)
		sysID = sysInfo.ID
	}
	return
}

func GetMailConf() (mail model.MailConf, sysID uint) {
	sysInfo, _ := dao2.NewSysConfDaoImpl().FindBySysType(model.SysMailConf)
	if sysInfo.ID > 0 {
		json.Unmarshal([]byte(sysInfo.Info), &mail)
		sysID = sysInfo.ID
	}
	return
}
