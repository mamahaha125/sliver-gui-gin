package dao

import (
	"github.com/jinzhu/gorm"
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
)

type SiteConfigDao interface {
	Insert(site model.PearConfig) error
	FindOne(configType model.PearConfigType) (*model.PearConfig, error)
	Update(site model.PearConfig, attr map[string]interface{}) error
}

func NewSiteConfigDaoImpl() SiteConfigDao {
	f := new(SiteConfigDaoImpl)
	return f
}

type SiteConfigDaoImpl struct {
}

func (f SiteConfigDaoImpl) Insert(site model.PearConfig) error {
	err := db.Instance().Model(model.PearConfig{}).Where("config_type = ?", site.ConfigType).FirstOrCreate(&site).Error
	return err
}

func (f SiteConfigDaoImpl) FindOne(configType model.PearConfigType) (*model.PearConfig, error) {
	var site model.PearConfig
	err := db.Instance().Model(model.PearConfig{}).Where("config_status = 1 AND config_type = ?", configType).First(&site).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &site, nil
}

func (f SiteConfigDaoImpl) Update(site model.PearConfig, attr map[string]interface{}) error {
	err := db.Instance().Model(&site).Where("id = ?", site.ID).Updates(attr).Error
	return err
}
