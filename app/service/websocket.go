// @Title       websocket.go
// @Description //TODO
// @Author      mamahaha125
// @Data        2023/8/29 22:52
package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/client/core"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"pear-admin-go/app/core/cmd"
	rpcs "pear-admin-go/app/core/rpc"
	"sync"
)

// @Update mamahaha125 2023-08-29 22:53:53
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 读写锁
var (
	rwLocker sync.RWMutex
)

type ShellWebsocket struct {
	conn    *websocket.Conn
	data    chan []byte
	session *clientpb.Session
	cmd     *cobra.Command
	con     *console.SliverConsoleClient
	nopty   bool
}

// @Title WebSocket
// @Description TODO
// @Author mamahaha125 2023-08-30 00:00:48
// @Update mamahaha125 2023-08-30 00:00:48
// @Param c
func WebSocket(writer http.ResponseWriter, request *http.Request) {
	isvalida := true //chektoken() 待
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	con := rpcs.GetInstance().GetCon()
	cmd.Cmd.ConsoleCmd()
	cmds := cmd.GetInstance().GetCmd()

	Conn := &ShellWebsocket{
		conn:    conn,
		data:    make(chan []byte, 10),
		session: con.ActiveTarget.GetSession(),
		cmd:     cmds,
		con:     con,
		nopty:   false,
	}

	//5.完成发送逻辑
	go sendProc(Conn)
	//6.完成接收逻辑
	go recvProc(Conn)

}

type filterReader struct {
	reader io.Reader
}

func newFilterReader(reader io.Reader) *filterReader {
	return &filterReader{reader}
}

func (r *filterReader) Read(data []byte) (int, error) {
	n, err := r.reader.Read(data)
	if err != nil {
		return n, err
	}

	// Remove Windows new line
	if n >= 2 {
		if data[n-2] == byte('\r') {
			data = data[0 : n-2]
			data = append(data, byte('\n'))
			n -= 1
		}
	}

	return n, nil
}

func sendProc(Conn *ShellWebsocket) {

	ctxTunnel, cancelTunnel := context.WithCancel(context.Background())
	rpcTunnel, err := Conn.con.Rpc.CreateTunnel(ctxTunnel, &sliverpb.Tunnel{
		SessionID: Conn.session.ID,
	})
	defer cancelTunnel()
	if err != nil {
		fmt.Println("%s\n", err, rpcTunnel)
		return
	}

	tunnel := core.GetTunnels().Start(rpcTunnel.TunnelID, rpcTunnel.SessionID)

	shell, err := Conn.con.Rpc.Shell(context.Background(), &sliverpb.ShellReq{
		Request:   Conn.con.ActiveTarget.Request(Conn.cmd),
		Path:      "",
		EnablePTY: !Conn.nopty,
		TunnelID:  tunnel.ID,
	})
	if err != nil {
		fmt.Println("%s\n", err)
		return
	}

	if shell.Response != nil && shell.Response.Err != "" {
		Conn.con.PrintErrorf("Error: %s\n", shell.Response.Err)
		_, err = Conn.con.Rpc.CloseTunnel(context.Background(), &sliverpb.Tunnel{
			TunnelID:  tunnel.ID,
			SessionID: Conn.session.ID,
		})
		if err != nil {
			fmt.Println("RPC Error: %s\n", err)
		}
		return
	}

	//defer tunnel.Close()
	fmt.Printf("Bound remote shell pid %d to tunnel %d", shell.Pid, shell.TunnelID)
	Conn.con.PrintInfof("Started remote shell with pid %d\n\n", shell.Pid)

	var outputBuffer bytes.Buffer
	fmt.Printf("Starting stdin/stdout shell ...")
	go func() {
		n, err := io.Copy(&outputBuffer, tunnel)
		fmt.Printf("Wrote %d bytes to stdout", n)
		if err != nil {
			fmt.Println("Error writing to stdout: %s", err)
			return
		}
	}()
	fmt.Printf("Reading from stdin ...")

	for {
		select {
		case data := <-Conn.data:
			fmt.Println("[ws] sendMsg >>>> msg:", string(data))
			reader := bytes.NewReader(data)
			n, err := io.Copy(tunnel, newFilterReader(reader))
			fmt.Printf("Read %d bytes from stdin", n)
			if err != nil && err != io.EOF {
				fmt.Println("Error reading from stdin: %s\n", err)
			}

			bufio.NewWriter(&outputBuffer).Flush()

			err = Conn.conn.WriteMessage(websocket.TextMessage, outputBuffer.Bytes())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	//// Create an RPC tunnel, then start it before binding the shell to the newly created tunnel
	//ctxTunnel, cancelTunnel := context.WithCancel(context.Background())
	//
	//rpcTunnel, err := con.Rpc.CreateTunnel(ctxTunnel, &sliverpb.Tunnel{
	//	SessionID: session.ID,
	//})
	//defer cancelTunnel()
	//noPty := false
	//// Start() takes an RPC tunnel and creates a local Reader/Writer tunnel object
	//tunnel := core.GetTunnels().Start(rpcTunnel.TunnelID, rpcTunnel.SessionID)
	//shellPath := ""
	//shell, err := con.Rpc.Shell(context.Background(), &sliverpb.ShellReq{
	//	Request:   con.ActiveTarget.Request(cmd),
	//	Path:      shellPath,
	//	EnablePTY: !noPty,
	//	TunnelID:  tunnel.ID,
	//})
	//defer tunnel.Close()
	//if shell.Response != nil && shell.Response.Err != "" {
	//	con.PrintErrorf("Error: %s\n", shell.Response.Err)
	//	_, err = con.Rpc.CloseTunnel(context.Background(), &sliverpb.Tunnel{
	//		TunnelID:  tunnel.ID,
	//		SessionID: session.ID,
	//	})
	//	if err != nil {
	//		con.PrintErrorf("RPC Error: %s\n", err)
	//	}
	//	return
	//}
	//var oldState *term.State
	//if !noPty {
	//	oldState, err = term.MakeRaw(0)
	//	fmt.Printf("Saving terminal state: %v", oldState)
	//	if err != nil {
	//		con.PrintErrorf("Failed to save terminal state")
	//		return
	//	}
	//}
	//n, err := io.Copy(tunnel, newFilterReader(os.Stdin))
	//term.Restore(0, oldState)
	//fmt.Println("Saving terminal state: %v", oldState)
	//var outputBuffer bytes.Buffer
	//go func() {
	//
	//	_, err := io.Copy(&outputBuffer, tunnel)
	//
	//	if err != nil {
	//		fmt.Println("Error writing to stdout: %s", err)
	//		return
	//	}
	//}()
	//
	//bufio.NewWriter(&outputBuffer).Flush()
	//
	//errs := Conn.WriteMessage(websocket.TextMessage, outputBuffer.Bytes())
	//if errs != nil {
	//	fmt.Println(errs)
	//	return
	//}
}

func recvProc(Conn *ShellWebsocket) {

	for {
		_, data, err := Conn.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		Conn.data <- data
		fmt.Println("[ws] <<<<", string(data))
	}
}
