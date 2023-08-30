package service

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"os"
	"path/filepath"
	rpcs "pear-admin-go/app/core/rpc"
)

func DeleteImplant(names []string) (string, error) {
	var err error
	con := rpcs.GetInstance().GetCon()
	builds, err := con.Rpc.ImplantBuilds(context.Background(), &commonpb.Empty{})
	if err != nil {
		return "", err
	}
	sliverNames := make(map[string]bool)
	for sliverName, _ := range builds.Configs {
		sliverNames[sliverName] = true
	}

	for _, name := range names {
		// Check if the sliver name exists
		if sliverNames[name] {
			_, err := con.Rpc.DeleteImplantBuild(context.Background(), &clientpb.DeleteReq{
				Name: name,
			})
			if err != nil {
				return "", err
			}

			filePath, _ := filepath.Abs(filepath.Join(".", "implants", name))
			err = os.Remove(filePath)
			if err != nil {
				// Accumulate file deletion errors
				fmt.Println(err)
			}
		}
	}
	return "name", err
}
