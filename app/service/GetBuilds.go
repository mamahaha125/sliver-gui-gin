package service

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"strconv"
)

type ImplantDetail struct {
	//Id     int    `json:"id" form:"id" zh:"ID"`
	Name   string   `json:"name,omitempty" form:"server_name" zh:"服务器名称"`
	Type   string   `json:"type,omitempty" form:"server_account"  zh:"用户名"`
	OSARCH string   `json:"osarch" form:"server_password" zh:"密码"`
	Format string   `json:"format,omitempty" form:"server_ip" zh:"服务器IP"`
	C2     []string `json:"c2,omitempty" form:"port"  zh:"端口号"`
	Debug  string   `json:"debug,omitempty" form:"private_key_src" zh:"私钥地址"` // ssh-keygen -t rsa -f pp_rsa

}

type SessionsDetail struct {
	//Id     int    `json:"id" form:"id" zh:"ID"`
	Name      string `json:"name,omitempty" form:"server_name" zh:"服务器名称"`
	Transport string `json:"transport,omitempty" form:"server_account"  zh:"用户名"`
	OSARCH    string `json:"osarch" form:"server_password" zh:"密码"`
	Hostname  string `json:"hostname,omitempty" form:"server_ip" zh:"服务器IP"`
	C2        string `json:"c2,omitempty" form:"port"  zh:"端口号"`
	Username  string `json:"username,omitempty" form:"private_key_src" zh:"私钥地址"` // ssh-keygen -t rsa -f pp_rsa
	Locale    string `json:"locale,omitempty" form:"port"  zh:"端口号"`
	Lastmsg   string `json:"lastmsg,omitempty" form:"port"  zh:"端口号"`
	ID        string `json:"id"`
}

func GetBuilds(con *console.SliverConsoleClient) ([]*ImplantDetail, int) {
	implants, err := con.Rpc.ImplantBuilds(context.Background(), &commonpb.Empty{})
	if err != nil {
		fmt.Println(err)
	}

	var Lists []*ImplantDetail
	for k, v := range implants.Configs {
		fmt.Println(v.Format.String())
		temp := &ImplantDetail{
			Name:   k,
			Type:   "session",
			OSARCH: v.GOOS + "/" + v.GOARCH,
			Format: v.Format.String(),
			C2: func() []string {
				var c2URLs []string
				//var c2 clientpb.ImplantC2

				for _, c2 := range v.C2 {

					c2URLs = append(c2URLs, c2.URL)
				}
				return c2URLs
			}(),
			Debug: strconv.FormatBool(v.Debug),
		}
		Lists = append(Lists, temp)
	}

	return Lists, len(implants.Configs)
}

func ResSessions(data []*clientpb.Session) ([]*SessionsDetail, int) {
	var Lists []*SessionsDetail
	for _, v := range data {
		temp := &SessionsDetail{
			Name:      v.Name,
			Transport: v.Transport,
			C2:        v.RemoteAddress,
			OSARCH:    v.OS,
			Username:  v.Username,
			Hostname:  v.Hostname,
			Lastmsg:   "shijian",
			Locale:    v.Locale,
			ID:        v.ID,
		}
		Lists = append(Lists, temp)
	}

	return Lists, len(Lists)
}
