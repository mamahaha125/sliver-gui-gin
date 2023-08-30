package main

import (
	"fmt"
	"sync"
)

type singleton struct {
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

func main() {
	var ss struct {
		s string `json:"s"`
		d string `json:"d"`
	}

	ss.s = "s"

	fmt.Println(ss.d)

}

//
//import (
//	"context"
//	"fmt"
//	"github.com/bishopfox/sliver/client/assets"
//	"github.com/bishopfox/sliver/client/console"
//	"github.com/bishopfox/sliver/client/transport"
//	"github.com/bishopfox/sliver/protobuf/commonpb"
//	"github.com/bishopfox/sliver/protobuf/rpcpb"
//	"google.golang.org/grpc"
//	"strconv"
//)
//
////// 查询本地客户端cfg文件
////func SearchConf() {
////	conf := GuiModels.UserBasic{}
////	ClientCfg := assets.GetConfigs()
////	conf.ConfigDir = assets.GetConfigDir()
////
////	for k := range ClientCfg {
////		conf.ConnectRpc = append(conf.ConnectRpc, k)
////	}
////	sort.Strings(conf.ConnectRpc)
////
////}
////
////// 选择文件进入rpc连接
////func ChooseConf() {
////	var rpc rpcpb.SliverRPCClient
////	var ln *grpc.ClientConn
////	var err error
////
////	configs := assets.GetConfigs()
////	clientName := "ll@192.168.1.101 (6f22d386f29b27df)"
////	ClientKey := configs[clientName]
////	fmt.Println(ClientKey.LPort)
////	rpc, ln, err = transport.MTLSConnect(ClientKey)
////	fmt.Println(rpc, "----", ln, "----", err)
////}
////
////// 初始化RPC
////func ConnectRPC() {
////	var rpc rpcpb.SliverRPCClient
////	var ln *grpc.ClientConn
////	var err error
////	//var cmd *cobra.Command
////
////	configs := assets.GetConfigs()
////	clientName := "new@192.168.200.101 (6b48954e9c0de3e1)"
////	ClientKey := configs[clientName]
////	fmt.Println("[用户配置端口]>>>", ClientKey.LPort)
////	rpc, ln, err = transport.MTLSConnect(ClientKey)
////	fmt.Println("[测试]>>>", rpc, "----", ln, "----", err)
////
////	con := console.NewConsole(false)
////	con.Rpc = rpc
////	//console.StartClient(con, rpc, command.ServerCommands(con, nil), command.SliverCommands(con), true)
////	jobs, errs := con.Rpc.GetJobs(context.Background(), &commonpb.Empty{})
////	if err != nil {
////		fmt.Println(errs)
////	}
////	fmt.Println("[jobs]>>>", jobs)
////
////}
//
//func GetJobs(rpcpb.SliverRPCClient) {
//	var rpc rpcpb.SliverRPCClient
//	var ln *grpc.ClientConn
//	var err error
//	//var cmd *cobra.Command
//
//	configs := assets.GetConfigs()
//	ClientKey := configs["new@192.168.200.101 (6b48954e9c0de3e1)"]
//	fmt.Println("[用户配置端口]>>>", ClientKey.LPort)
//	rpc, ln, err = transport.MTLSConnect(ClientKey)
//	fmt.Println("[测试]>>>", rpc, "----", ln, "----", err)
//
//	con := console.NewConsole(false)
//	con.Rpc = rpc
//	//console.StartClient(con, rpc, command.ServerCommands(con, nil), command.SliverCommands(con), true)
//	jobs, errs := con.Rpc.GetJobs(context.Background(), &commonpb.Empty{})
//	if err != nil {
//		fmt.Println(errs)
//	}
//	fmt.Println("[jobs]>>>", jobs)
//}
//
//func GetSessions(rpcpb.SliverRPCClient) {
//	var rpc rpcpb.SliverRPCClient
//	var ln *grpc.ClientConn
//	var err error
//	//var cmd *cobra.Command
//
//	configs := assets.GetConfigs()
//	ClientKey := configs["new@192.168.200.101 (6b48954e9c0de3e1)"]
//	fmt.Println("[用户配置端口]>>>", ClientKey.LPort)
//	rpc, ln, err = transport.MTLSConnect(ClientKey)
//	fmt.Println("[测试]>>>", rpc, "----", ln, "----", err)
//
//	con := console.NewConsole(false)
//	con.Rpc = rpc
//	//console.StartClient(con, rpc, command.ServerCommands(con, nil), command.SliverCommands(con), true)
//	jobs, errs := con.Rpc.GetSessions(context.Background(), &commonpb.Empty{})
//	if err != nil {
//		fmt.Println(errs)
//	}
//	if jobs.Sessions != nil {
//		fmt.Println("none")
//	}
//	fmt.Println("[sessinos]>>>", jobs.Sessions[0])
//}
//
//func GetBuilds(rpcpb.SliverRPCClient) {
//	var rpc rpcpb.SliverRPCClient
//	var ln *grpc.ClientConn
//	var err error
//	//var cmd *cobra.Command
//
//	configs := assets.GetConfigs()
//	ClientKey := configs["new@192.168.200.101 (6b48954e9c0de3e1)"]
//	fmt.Println("[用户配置端口]>>>", ClientKey.LPort)
//	rpc, ln, err = transport.MTLSConnect(ClientKey)
//	fmt.Println("[测试]>>>", rpc, "----", ln, "----", err)
//
//	con := console.NewConsole(false)
//	con.Rpc = rpc
//	//console.StartClient(con, rpc, command.ServerCommands(con, nil), command.SliverCommands(con), true)
//	jobs, errs := con.Rpc.ImplantBuilds(context.Background(), &commonpb.Empty{})
//	if err != nil {
//		fmt.Println(errs)
//	}
//
//	fmt.Println("[sessinos]>>>", jobs)
//}
//
//func GetWebsite(rpcpb.SliverRPCClient) {
//	var rpc rpcpb.SliverRPCClient
//	var ln *grpc.ClientConn
//	var err error
//	//var cmd *cobra.Command
//
//	configs := assets.GetConfigs()
//	ClientKey := configs["new@192.168.200.101 (6b48954e9c0de3e1)"]
//	fmt.Println("[用户配置端口]>>>", ClientKey.LPort)
//	rpc, ln, err = transport.MTLSConnect(ClientKey)
//	fmt.Println("[测试]>>>", rpc, "----", ln, "----", err)
//
//	con := console.NewConsole(false)
//	con.Rpc = rpc
//	//console.StartClient(con, rpc, command.ServerCommands(con, nil), command.SliverCommands(con), true)
//	jobs, errs := con.Rpc.Websites(context.Background(), &commonpb.Empty{})
//	if err != nil {
//		fmt.Println(errs)
//	}
//
//	fmt.Println("[sessinos]>>>", jobs)
//}
//
//func main() {
//	//GetSessions()
//
//	a, _ := strconv.ParseBool("")
//	//if a {
//	fmt.Println(a)
//	//}
//}
