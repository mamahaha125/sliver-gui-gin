package dao

import (
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
	"strings"
)

type AuthDao interface {
	Insert(auth model.Auth) (authID uint, err error)
	Find(page, pageSize int, filters ...interface{}) (auths []model.Auth, total int64)
	FindOne(id int) (auth model.Auth, err error)
	FindChildNode(id int) (int, error)
	Update(auth model.Auth, attr map[string]interface{}) error
	Delete(id int) error
	DeleteUse() error
	AuthList() ([]model.AuthListResp, error)
}

func NewAuthDaoImpl() AuthDao {
	auth := new(AuthDaoImpl)
	return auth
}

type AuthDaoImpl struct {
}

func (a *AuthDaoImpl) AuthList() (authResp []model.AuthListResp, err error) {
	client := db.Instance()
	err = client.Raw("SELECT a.id AS cate_id, a.auth_name AS cate_name, b.id AS menu_id, " +
		"b.auth_name AS menu_name, b.auth_url AS menu_url FROM auth a JOIN auth b ON b.pid = a.id " +
		"WHERE a.power_type = '0' AND b.power_type = 1 ORDER BY a.id ASC, a.sort ASC").Scan(&authResp).Error
	if err != nil {
		return nil, err
	}
	return authResp, nil
}

func (a *AuthDaoImpl) DeleteUse() error {
	client := db.Instance()
	var secondAuth []model.Auth
	client.Where("pid = 20").Find(&secondAuth)
	if len(secondAuth) > 0 {
		for _, v := range secondAuth {
			client.Where("pid = ?", v.ID).Delete(model.Auth{})
			client.Where("id = ?", v.ID).Delete(model.Auth{})
		}
	}

	return nil
}

func (a *AuthDaoImpl) Insert(auth model.Auth) (authID uint, err error) {
	client := db.Instance()
	err = client.Create(&auth).Error
	return auth.ID, nil
}

func (a *AuthDaoImpl) Find(page, pageSize int, filters ...interface{}) (auths []model.Auth, total int64) {
	offset := (page - 1) * pageSize
	client := db.Instance()
	client = client.Model(model.Auth{})
	var queryArr []string
	var values []interface{}
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			queryArr = append(queryArr, gconv.String(filters[k]))
			values = append(values, filters[k+1])
		}
	}
	client.Model(model.Auth{}).Where(strings.Join(queryArr, " AND "), values...).Order("power_type,sort,pid").Limit(pageSize).Offset(offset).Find(&auths)
	client.Model(model.Auth{}).Where(strings.Join(queryArr, " AND "), values...).Count(&total)
	return
}

func (a *AuthDaoImpl) FindOne(id int) (auth model.Auth, err error) {
	client := db.Instance()
	client.First(&auth, id)
	return auth, client.Error
}

func (a *AuthDaoImpl) FindChildNode(id int) (int, error) {
	var count int
	client := db.Instance()
	client.Model(model.Auth{}).Where("status = 1 AND pid = ?", id).Count(&count)
	return count, client.Error
}

func (a *AuthDaoImpl) Update(auth model.Auth, attr map[string]interface{}) error {
	client := db.Instance()
	if _, ok := attr["pid"]; ok {
		attr["pid"] = gconv.Int(attr["pid"])
	}
	if _, ok := attr["sort"]; ok {
		attr["sort"] = gconv.Int(attr["sort"])
	}
	if _, ok := attr["power_type"]; ok {
		attr["power_type"] = gconv.Int(attr["power_type"])
	}
	if _, ok := attr["is_show"]; ok {
		attr["is_show"] = gconv.Int(attr["is_show"])
	}
	client.Model(&auth).Omit("id").Updates(attr)
	return client.Error
}

func (a *AuthDaoImpl) Delete(id int) error {
	client := db.Instance()
	client.Where("id = ?", id).Delete(model.Auth{})
	return client.Error
}
