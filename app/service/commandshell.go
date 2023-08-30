// @Title       commandshell.go
// @Description //TODO
// @Author      mamahaha125
// @Data        2023/8/29 12:23
package service

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	"pear-admin-go/app/core/cmd"
	rpcs "pear-admin-go/app/core/rpc"
)

func ExecCommand(args []string) (string, error) {
	var path []string
	con := rpcs.GetInstance().GetCon()
	cmds := cmd.Cmd.GetCmd()

	var err error
	var exec *sliverpb.Execute

	if len(args) > 1 {
		path = []string{args[1]}
	} else {
		path = []string{}
	}

	exec, err = con.Rpc.Execute(context.Background(), &sliverpb.ExecuteReq{
		Request: con.ActiveTarget.Request(cmds),
		Path:    args[0],
		Args:    path,
		Output:  true,
		Stderr:  "",
		Stdout:  "",
	})
	if err != nil {
		fmt.Println()
	}

	return string(exec.Stdout), nil
}
