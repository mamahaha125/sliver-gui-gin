package rpc

import (
	"context"
	"fmt"
	sliverconsole "github.com/bishopfox/sliver/client/console"
	consts "github.com/bishopfox/sliver/client/constants"
	"github.com/bishopfox/sliver/client/core"
	"github.com/bishopfox/sliver/client/prelude"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"google.golang.org/protobuf/proto"
	"io"
	"strings"
	"sync"
	"time"
)

type RpcController struct {
	con *sliverconsole.SliverConsoleClient
}

var (
	Con  *RpcController
	once sync.Once
)

func GetInstance() *RpcController {
	once.Do(func() {
		Con = &RpcController{}
	})
	return Con
}

func (Con *RpcController) GetCon() *sliverconsole.SliverConsoleClient {
	return Con.con
}

func (Con *RpcController) InitCon(con *sliverconsole.SliverConsoleClient) {
	GetInstance().con = con
}

func (Con *RpcController) InitRpc(rpc rpcpb.SliverRPCClient) {
	GetInstance().con.Rpc = rpc
}

func (Con *RpcController) StartEventLoop() {
	eventStream, err := Con.con.Rpc.Events(context.Background(), &commonpb.Empty{})
	if err != nil {
		fmt.Printf(sliverconsole.Warn+"%s\n", err)
		return
	}
	for {
		event, err := eventStream.Recv()
		if err == io.EOF || event == nil {
			return
		}

		go Con.triggerEventListeners(event)

		// Trigger event based on type
		switch event.EventType {

		case consts.CanaryEvent:
			Con.con.PrintEventErrorf(sliverconsole.Bold+"WARNING: %s%s has been burned (DNS Canary)", sliverconsole.Normal, event.Session.Name)
			sessions := Con.con.GetSessionsByName(event.Session.Name)
			for _, session := range sessions {
				shortID := strings.Split(session.ID, "-")[0]
				Con.con.PrintErrorf("\t🔥 Session %s is affected", shortID)
			}

		case consts.WatchtowerEvent:
			msg := string(event.Data)
			Con.con.PrintEventErrorf(sliverconsole.Bold+"WARNING: %s%s has been burned (seen on %s)", sliverconsole.Normal, event.Session.Name, msg)
			sessions := Con.con.GetSessionsByName(event.Session.Name)
			for _, session := range sessions {
				shortID := strings.Split(session.ID, "-")[0]
				Con.con.PrintErrorf("\t🔥 Session %s is affected", shortID)
			}

		case consts.JoinedEvent:
			if Con.con.Settings.UserConnect {
				Con.con.PrintInfof("%s has joined the game", event.Client.Operator.Name)
			}
		case consts.LeftEvent:
			if Con.con.Settings.UserConnect {
				Con.con.PrintInfof("%s left the game", event.Client.Operator.Name)
			}

		case consts.JobStoppedEvent:
			job := event.Job
			Con.con.PrintErrorf("Job #%d stopped (%s/%s)", job.ID, job.Protocol, job.Name)

		case consts.SessionOpenedEvent:
			session := event.Session
			currentTime := time.Now().Format(time.RFC1123)
			shortID := strings.Split(session.ID, "-")[0]
			Con.con.PrintEventInfof("Session %s %s - %s (%s) - %s/%s - %v",
				shortID, session.Name, session.RemoteAddress, session.Hostname, session.OS, session.Arch, currentTime)

			// Prelude Operator
			if prelude.ImplantMapper != nil {
				err = prelude.ImplantMapper.AddImplant(session, nil)
				if err != nil {
					Con.con.PrintErrorf("Could not add session to Operator: %s", err)
				}
			}

		case consts.SessionUpdateEvent:
			session := event.Session
			currentTime := time.Now().Format(time.RFC1123)
			shortID := strings.Split(session.ID, "-")[0]
			Con.con.PrintInfof("Session %s has been updated - %v", shortID, currentTime)

		case consts.SessionClosedEvent:
			session := event.Session
			currentTime := time.Now().Format(time.RFC1123)
			shortID := strings.Split(session.ID, "-")[0]
			Con.con.PrintEventErrorf("Lost session %s %s - %s (%s) - %s/%s - %v",
				shortID, session.Name, session.RemoteAddress, session.Hostname, session.OS, session.Arch, currentTime)
			activeSession := Con.con.ActiveTarget.GetSession()
			core.GetTunnels().CloseForSession(session.ID)
			core.CloseCursedProcesses(session.ID)
			if activeSession != nil && activeSession.ID == session.ID {
				Con.con.ActiveTarget.Set(nil, nil)
				Con.con.PrintErrorf("Active session disconnected")
			}
			if prelude.ImplantMapper != nil {
				err = prelude.ImplantMapper.RemoveImplant(session)
				if err != nil {
					Con.con.PrintErrorf("Could not remove session from Operator: %s", err)
				}
				Con.con.PrintInfof("Removed session %s from Operator", session.Name)
			}

		case consts.BeaconRegisteredEvent:
			beacon := &clientpb.Beacon{}
			proto.Unmarshal(event.Data, beacon)
			currentTime := time.Now().Format(time.RFC1123)
			shortID := strings.Split(beacon.ID, "-")[0]
			Con.con.PrintEventInfof("Beacon %s %s - %s (%s) - %s/%s - %v",
				shortID, beacon.Name, beacon.RemoteAddress, beacon.Hostname, beacon.OS, beacon.Arch, currentTime)

			// Prelude Operator
			if prelude.ImplantMapper != nil {
				err = prelude.ImplantMapper.AddImplant(beacon, func(taskID string, cb func(*clientpb.BeaconTask)) {
					Con.con.AddBeaconCallback(taskID, cb)
				})
				if err != nil {
					Con.con.PrintErrorf("Could not add beacon to Operator: %s", err)
				}
			}

		case consts.BeaconTaskResultEvent:
			Con.triggerBeaconTaskCallback(event.Data)

		}

		Con.triggerReactions(event)
	}
}

func (Con *RpcController) triggerEventListeners(event *clientpb.Event) {
	Con.con.EventListeners.Range(func(key, value interface{}) bool {
		listener := value.(chan *clientpb.Event)
		listener <- event // Do not block while sending the event to the listener
		return true
	})
}

// triggerBeaconTaskCallback - Triggers the callback for a beacon task
func (Con *RpcController) triggerBeaconTaskCallback(data []byte) {
	task := &clientpb.BeaconTask{}
	err := proto.Unmarshal(data, task)
	if err != nil {
		Con.con.PrintErrorf("\rCould not unmarshal beacon task: %s", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	beacon, _ := Con.con.Rpc.GetBeacon(ctx, &clientpb.Beacon{ID: task.BeaconID})

	// If the callback is not in our map then we don't do anything, the beacon task
	// was either issued by another operator in multiplayer mode or the client process
	// was restarted between the time the task was created and when the server got the result
	Con.con.BeaconTaskCallbacksMutex.Lock()
	defer Con.con.BeaconTaskCallbacksMutex.Unlock()
	if callback, ok := Con.con.BeaconTaskCallbacks[task.ID]; ok {
		if Con.con.Settings.BeaconAutoResults {
			if beacon != nil {
				Con.con.PrintEventSuccessf("%s completed task %s", beacon.Name, strings.Split(task.ID, "-")[0])
			}
			task_content, err := Con.con.Rpc.GetBeaconTaskContent(ctx, &clientpb.BeaconTask{
				ID: task.ID,
			})
			Con.con.Printf(sliverconsole.Clearln + "\r")
			if err == nil {
				callback(task_content)
			} else {
				Con.con.PrintErrorf("Could not get beacon task content: %s", err)
			}
			Con.con.Println()
		}
		delete(Con.con.BeaconTaskCallbacks, task.ID)
	}
}

func (Con *RpcController) triggerReactions(event *clientpb.Event) {
	reactions := core.Reactions.On(event.EventType)
	if len(reactions) == 0 {
		return
	}

	// We need some special handling for SessionOpenedEvent to
	// set the new session as the active session
	currentActiveSession, currentActiveBeacon := Con.con.ActiveTarget.Get()
	defer func() {
		Con.con.ActiveTarget.Set(currentActiveSession, currentActiveBeacon)
	}()

	if event.EventType == consts.SessionOpenedEvent {
		Con.con.ActiveTarget.Set(nil, nil)

		Con.con.ActiveTarget.Set(event.Session, nil)
	} else if event.EventType == consts.BeaconRegisteredEvent {
		Con.con.ActiveTarget.Set(nil, nil)

		beacon := &clientpb.Beacon{}
		proto.Unmarshal(event.Data, beacon)
		Con.con.ActiveTarget.Set(nil, beacon)
	}

	for _, reaction := range reactions {
		for _, line := range reaction.Commands {
			Con.con.PrintInfof(sliverconsole.Bold+"Execute reaction: '%s'"+sliverconsole.Normal, line)
			err := Con.con.App.ActiveMenu().RunCommand(line)
			if err != nil {
				Con.con.PrintErrorf("Reaction command error: %s\n", err)
			}
		}
	}
}
