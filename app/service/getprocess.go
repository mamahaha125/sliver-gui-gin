// @Title       getprocess.go
// @Description //TODO
// @Author      mamahaha125
// @Data        2023/8/27 15:58
package service

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	"pear-admin-go/app/core/cmd"
	rpcs "pear-admin-go/app/core/rpc"
)

type PsList struct {
	//Title string `json:"title"`
	//Type    string `json:"type"`
	//Mode    string `json:"mode"`
	Title string `json:"title"`
	Pid   int32  `json:"pid"`
	Ppid  int32  `json:"ppid"`

	Owner   string `json:"owner"`
	Arch    string `json:"arch"`
	CmdLine string `json:"cmdLine"`

	ModTime int64  `json:"modTime"`
	Size    int64  `json:"size"`
	Path    string `json:"path"`
	Id      string `json:"id"`
}

// @Title PSList
// @Description 获取进程列表
// @Author mamahaha125 2023-08-28 21:39:03
// @Update mamahaha125 2023-08-28 21:39:03
// @Return []PsList
func PSList() []PsList {
	//clientName := "new@192.168.200.101 (6b48954e9c0de3e1)"
	//
	//rpc, cons, _ := ConnectRpc(clientName)
	//SetCon(cons, rpc)
	con := rpcs.GetInstance().GetCon()
	cmds := cmd.GetInstance().GetCmd()

	ps, err := con.Rpc.Ps(context.Background(), &sliverpb.PsReq{
		Request: con.ActiveTarget.Request(cmds),
	})
	if err != nil {
		fmt.Println(err)
	}

	var List []PsList
	var PS PsList
	var Cmd string
	var num int
	for _, v := range ps.Processes {

		num = num + 1
		PS.Id = fmt.Sprintf("%03d", num)
		PS.Title = v.Executable
		PS.Pid = v.Pid
		PS.Ppid = v.Ppid
		PS.Owner = v.Owner
		PS.Arch = v.Architecture
		//for name := range v.CmdLine {
		//	Cmd = fmt.Sprintf("%s", name) + " " + Cmd
		//}
		PS.CmdLine = Cmd
		List = append(List, PS)
	}

	return List
}
