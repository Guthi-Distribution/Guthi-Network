package nodes

import (
	"fmt"
	"net"
	"os"
)

// So each node should work as server and client simultaneously

type NetworkNode struct {
	NodeID int
	Name   string
	Socket *net.TCPAddr // TCP address of the current node
}

// Test code
func CheckTCPConnection() {
	tcp_info := "localhost:8080"
	tcp_addr, err := net.ResolveTCPAddr("tcp", tcp_info)
	if err != nil {
		println("Cannot resolve TCP addr info : ", err)
		os.Exit(1)
	}

	socket, err := net.DialTCP("tcp", nil, tcp_addr)

	if err != nil {
		println("Failed to dial tcp server", err.Error())
		os.Exit(1)
	}
	_, err = socket.Write([]byte("Send me time bro"))

	time := make([]byte, 2048)
	// Reading from the server
	for {
		count, err := socket.Read(time)
		if err != nil {
			println("Failed to read from the server")
			break
		}
		println("Time received is : ", string(time[0:count]))
	}
	socket.Close()
}

func IntiateTCPConnection(node *NetworkNode) *net.TCPConn {
	tcp_con, err := net.DialTCP("tcp", nil, node.Socket)
	if err != nil {
		fmt.Println("Failed to initiate tcp connection with the host : ", node)
		return nil
	}
	return tcp_con
}

func ListenForTCPConnection(node *NetworkNode) {
	listener, err := net.ListenTCP("tcp", node.Socket)

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
	for {
		conn, _ := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Failed to Accept the incoming connection")
			break
		}
		go HandleTCPConnection(conn)
	}
	listener.Close()
}

func HandleTCPConnection(tcp_connection *net.TCPConn) {

}
