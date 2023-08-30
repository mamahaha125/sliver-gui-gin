package service

import (
	"encoding/json"
	"pear-admin-go/app/core/cache"
	"pear-admin-go/app/dao"
	e2 "pear-admin-go/app/global/e"
	"pear-admin-go/app/model"
)

func GetPearConfig() (*model.PearConfigForm, error) {
	pear, err := dao.NewSiteConfigDaoImpl().FindOne(model.PearSiteConfig)
	if err != nil {
		return nil, err
	}
	var data model.PearConfigForm
	err = json.Unmarshal([]byte(pear.ConfigData), &data)
	if err != nil {
		return nil, err
	}
	cache.Instance().Set(e2.PearConfigCache, data, 0)
	return &data, nil
}
