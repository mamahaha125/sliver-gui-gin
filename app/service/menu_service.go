package service

import (
	"github.com/gin-gonic/gin"
	dao2 "pear-admin-go/app/dao"
	"pear-admin-go/app/model"
)

type SelfMenuData struct {
	ParentData []map[string]interface{}
	ChildData  []map[string]interface{}
	AllowUrl   string
}

type CacheMenu struct {
	List1    []map[string]interface{}
	List2    []map[string]interface{}
	AllowUrl string
}

type CacheMenuV2 struct {
	MenuResp []model.MenuResp `json:"menu_resp"`
	AllowUrl string           `json:"allow_url"`
}

func MenuServiceV2(c *gin.Context) (cacheMenu CacheMenuV2) {
	///
	/// 修改部分
	filters := make([]interface{}, 0)
	filters = append(filters, "status = ?", 1)
	filters = append(filters, "power_type = ?", 0)
	auths, _ := dao2.NewAuthDaoImpl().Find(1, 5000, filters...)
	var menu []model.MenuResp
	for _, auth := range auths {
		cf := make([]interface{}, 0)
		cf = append(cf, "status = ?", 1)
		cf = append(cf, "power_type = ?", 1)
		cf = append(cf, "pid = ?", auth.ID)
		childs, _ := dao2.NewAuthDaoImpl().Find(1, 5000, cf...)
		var childsResp []model.MenuChildrenResp
		for _, child := range childs {
			childsResp = append(childsResp, model.MenuChildrenResp{
				ID:       int(child.ID),
				Title:    child.AuthName,
				Type:     child.PowerType,
				Icon:     child.Icon,
				Href:     child.AuthUrl,
				OpenType: "_iframe",
			})
		}
		menu = append(menu, model.MenuResp{
			ID:       int(auth.ID),
			Title:    auth.AuthName,
			Type:     auth.PowerType,
			Icon:     auth.Icon,
			Href:     auth.AuthUrl,
			Children: childsResp,
		})
	}

	cacheMenu.AllowUrl = ""
	cacheMenu.MenuResp = menu

	return cacheMenu
	///
	//user := GetProfile(c)
	//hasMenu, found := cache.Instance().Get(e2.MenuCache + gconv.String(user.ID))
	//if found && hasMenu != nil {
	//	cacheMenu = hasMenu.(CacheMenuV2)
	//} else {
	//	var authId []string
	//	result := GetAuth(user)
	//	allowUrl := ""
	//	for _, v := range result {
	//		if !pkg.IsContain([]string{"", "/"}, v.AuthUrl) {
	//			if allowUrl == "" {
	//				allowUrl += v.AuthUrl
	//			} else {
	//				allowUrl += "," + v.AuthUrl
	//			}
	//		}
	//		authId = append(authId, strconv.Itoa(int(v.ID)))
	//	}
	//	filters := make([]interface{}, 0)
	//	filters = append(filters, "status = ?", 1)
	//	filters = append(filters, "power_type = ?", 0)
	//	auths, _ := dao2.NewAuthDaoImpl().Find(1, 5000, filters...)
	//	var menu []model.MenuResp
	//	for _, auth := range auths {
	//		if IsAdmin(user) == false && pkg.IsContain(authId, strconv.Itoa(int(auth.ID))) == false { // 管理员不限制
	//			continue
	//		}
	//		cf := make([]interface{}, 0)
	//		cf = append(cf, "status = ?", 1)
	//		cf = append(cf, "power_type = ?", 1)
	//		cf = append(cf, "pid = ?", auth.ID)
	//		childs, _ := dao2.NewAuthDaoImpl().Find(1, 5000, cf...)
	//		var childsResp []model.MenuChildrenResp
	//		for _, child := range childs {
	//			if IsAdmin(user) == false && pkg.IsContain(authId, strconv.Itoa(int(child.ID))) == false { // 管理员不限制
	//				continue
	//			}
	//			childsResp = append(childsResp, model.MenuChildrenResp{
	//				ID:       int(child.ID),
	//				Title:    child.AuthName,
	//				Type:     child.PowerType,
	//				Icon:     child.Icon,
	//				Href:     child.AuthUrl,
	//				OpenType: "_iframe",
	//			})
	//		}
	//		menu = append(menu, model.MenuResp{
	//			ID:       int(auth.ID),
	//			Title:    auth.AuthName,
	//			Type:     auth.PowerType,
	//			Icon:     auth.Icon,
	//			Href:     auth.AuthUrl,
	//			Children: childsResp,
	//		})
	//	}
	//
	//	cacheMenu.AllowUrl = allowUrl
	//	cacheMenu.MenuResp = menu
	//	cache.Instance().Set(e2.MenuCache+gconv.String(user.ID), cacheMenu, time.Hour)
	//}
	//return cacheMenu
}
