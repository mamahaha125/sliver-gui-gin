// @title       mtls_start.go
// @description //TODO
// @author      mamahaha125
// @data        2023/8/24 6:02
package service

import (
	"context"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	rpcs "pear-admin-go/app/core/rpc"
)

// @Title MtlsListen
// @Description 开启mtls监听
// @Author mamahaha125 2023-08-24 06:24:37
// @Update mamahaha125 2023-08-24 06:24:37
// @Param lhost
// @Param lport
// @Param persistent
// @Return uint32
// @Return error
func MtlsListen(lhost string, lport uint32, persistent bool) (uint32, error) {
	con := rpcs.GetInstance().GetCon()
	var errs error
	mtls, err := con.Rpc.StartMTLSListener(context.Background(), &clientpb.MTLSListenerReq{
		Host:       lhost,
		Port:       lport,
		Persistent: persistent,
	})
	var nils uint32
	if err != nil {
		return nils, errs
	}
	return mtls.JobID, nil
}
