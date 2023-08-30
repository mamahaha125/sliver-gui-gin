package service

import (
	"fmt"
	"github.com/bishopfox/sliver/client/assets"
	"github.com/bishopfox/sliver/client/command"
	sliverconsole "github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"log"
	"os"
	"path/filepath"
	"pear-admin-go/app/core/cmd"
	rpcs "pear-admin-go/app/core/rpc"
	"time"
)

// @Title SetCon
// @Description 开启sliver rpc，cmd连接与监听
// @Author mamahaha125 2023-08-29 16:34:20
// @Update mamahaha125 2023-08-29 16:34:20
// @Param con
// @Param rpc
func SetCon(con *sliverconsole.SliverConsoleClient, rpc rpcpb.SliverRPCClient) {
	Con := rpcs.GetInstance()
	Con.InitCon(con)
	Con.InitRpc(rpc)
	cmd.Cmd.ConsoleCmd()
	sliverconsole.StartClient(Con.GetCon(), rpc, command.ServerCommands(Con.GetCon(), nil), command.SliverCommands(Con.GetCon()), false)

	//go Con.StartEventLoop()
	//go core.TunnelLoop(rpc)

}

func getConsoleLogFile() *os.File {
	logsDir := assets.GetConsoleLogsDir()
	dateTime := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("%s.log", dateTime))
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Could not open log file: %s", err)
	}
	return logFile
}

func getConsoleAsciicastFile() *os.File {
	logsDir := assets.GetConsoleLogsDir()
	dateTime := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("asciicast_%s.log", dateTime))
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Could not open log file: %s", err)
	}
	return logFile
}
