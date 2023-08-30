package service

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/protobuf/clientpb"
)

func MtlsListener(con *console.SliverConsoleClient) (*clientpb.MTLSListener, error) {
	var lhost string
	var lport uint32
	var persistent bool
	mtls, err := con.Rpc.StartMTLSListener(context.Background(), &clientpb.MTLSListenerReq{
		Host:       lhost,
		Port:       lport,
		Persistent: persistent,
	})
	if err != nil {
		fmt.Println(err)
	}
	return mtls, nil
}
