package service

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
)

func GetJobs(con *console.SliverConsoleClient) *clientpb.Jobs {
	jobs, err := con.Rpc.GetJobs(context.Background(), &commonpb.Empty{})
	if err != nil {
		return nil
	}
	return jobs
}

func GetSessions(con *console.SliverConsoleClient) ([]*clientpb.Session, error) {
	jobs, err := con.Rpc.GetSessions(context.Background(), &commonpb.Empty{})
	var errs error
	if err != nil {
		return nil, errs
	}
	if len(jobs.Sessions) == 0 {
		return nil, errs
	}
	return jobs.Sessions, nil
}

func GetWebsite(con *console.SliverConsoleClient) *clientpb.Websites {
	jobs, err := con.Rpc.Websites(context.Background(), &commonpb.Empty{})
	if err != nil {
		fmt.Println(err)
	}

	return jobs
}

func GetCanaries() {

}
