package service

import (
	"fmt"
	"github.com/bishopfox/sliver/client/assets"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/client/transport"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"google.golang.org/grpc"
)

func ConnectRpc(clientName string) (rpcpb.SliverRPCClient, *console.SliverConsoleClient, *grpc.ClientConn) {
	var rpc rpcpb.SliverRPCClient
	var ln *grpc.ClientConn
	var err error

	configs := assets.GetConfigs()

	ClientKey := configs[clientName]
	rpc, ln, err = transport.MTLSConnect(ClientKey)
	con := console.NewConsole(false)
	con.Rpc = rpc
	fmt.Println(rpc, "----", ln, "----", err)
	return rpc, con, ln
}
