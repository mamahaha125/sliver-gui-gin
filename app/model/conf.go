package model

import (
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"google.golang.org/grpc"
)

type UserBasic struct {
	ConnectRpc []string               //cfg文件
	ConfigDir  string                 //配置文件地址
	Rpc        *rpcpb.SliverRPCClient //RPC链接
	link       *grpc.ClientConn       //RPC链接
}
