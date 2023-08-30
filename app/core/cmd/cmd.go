// @Title       cmd.go
// @Description //TODO
// @Author      mamahaha125
// @Data        2023/8/29 16:23
package cmd

import (
	"github.com/spf13/cobra"
	"sync"
)

type CmdConn struct {
	cmd *cobra.Command
}

var (
	Cmd  *CmdConn
	once sync.Once
)

// @Title consoleCmd
// @Description 初始化cobra
// @Author mamahaha125 2023-08-28 21:38:45
// @Update mamahaha125 2023-08-29 16:24:00
// @Param con
// @Return *cobra.Command
func (Cmd *CmdConn) ConsoleCmd() *CmdConn {
	return &CmdConn{&cobra.Command{
		Use:   "console",
		Short: "Start the sliver client console",
	}}
	//consoleCmd.RunE, consoleCmd.PersistentPostRunE = consoleRunnerCmd(con, true)
}

// @Title GetInstance
// @Description 单例
// @Author mamahaha125 2023-08-29 16:32:08
// @Update mamahaha125 2023-08-29 16:32:08
// @Return *CmdConn
func GetInstance() *CmdConn {
	once.Do(func() {
		Cmd = &CmdConn{}
	})
	return Cmd
}

// @Title GetCon
// @Description TODO
// @Author mamahaha125 2023-08-29 16:32:06
// @Update mamahaha125 2023-08-29 16:32:06
// @Return *cobra.Command
func (Cmd *CmdConn) GetCmd() *cobra.Command {
	return Cmd.cmd
}
