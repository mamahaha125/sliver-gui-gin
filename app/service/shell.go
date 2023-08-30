package service

import (
	"github.com/bishopfox/sliver/protobuf/clientpb"
	rpcs "pear-admin-go/app/core/rpc"
)

func ChooseSessionByID(name string) error {
	var errs error
	con := rpcs.GetInstance().GetCon()
	sessions, err := GetSessions(con)
	if err != nil {
		return errs
	}

	for _, session := range sessions {
		if session.ID == name {
			con.ActiveTarget.Set(session, nil)
			return nil
		}
	}
	return errs
}

func GetCurrentSession() *clientpb.Session {
	con := rpcs.GetInstance().GetCon()
	if session := con.ActiveTarget.GetSession(); session != nil {
		return session
	}
	return nil
}
