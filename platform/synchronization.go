package platform

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"time"
)

// TODO: Connection timeout feature complete??
var pending_connection_time = make(map[string]uint64)

// TODO: Need to add timestamp
type EchoMessage struct {
	AddrFrom     string
	ConnectionId uint64
}

type EchoReply struct {
	AddrFrom     string
	ConnectionId uint64
	IsReply      bool // to indicate either the reply is just a reply or reply to a reply
}

func SendEchoMessage(addr string, net_platform *NetworkPlatform) error {
	rand_num, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		return err
	}

	pending_connection_time[addr] = uint64(time.Now().Unix())
	payload := EchoMessage{
		AddrFrom:     net_platform.GetNodeAddress(),
		ConnectionId: rand_num.Uint64(),
	}
	data := GobEncode(payload)
	data = append(CommandStringToBytes("echo"), data...)
	return sendDataToAddress(addr, data, net_platform)
}

func HandleEchoMessage(request []byte, net_platform *NetworkPlatform) error {
	var recv_payload EchoMessage
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&recv_payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}
	// if the receiving address is the self address, then it is send
	send_payload := EchoReply{
		AddrFrom:     net_platform.GetNodeAddress(),
		ConnectionId: recv_payload.ConnectionId + 1,
		IsReply:      false,
	}

	return sendDataToAddress(recv_payload.AddrFrom, append(CommandStringToBytes("echo_reply"), GobEncode(send_payload)...), net_platform)
}

func HandleEchoReply(request []byte, net_platform *NetworkPlatform) error {
	var payload EchoReply
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)

	if !payload.IsReply {
		// echo reply is recieved
		send_payload := EchoReply{
			AddrFrom: net_platform.GetNodeAddress(),
			IsReply:  true,
		}
		pending_connection_time[payload.AddrFrom] = uint64(time.Now().Unix())
		err := sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("echo_reply"), GobEncode(send_payload)...), net_platform)
		if err != nil {
			return err
		}
	}

	// if the data is received delete it from the pending connection information
	if _, id := pending_connection_time[payload.AddrFrom]; id {
		delete(pending_connection_time, payload.AddrFrom)
	}
	return nil
}

func CheckForResponse(net_platform *NetworkPlatform) {
	for node, send_time := range pending_connection_time {
		curr_time := uint64(time.Now().Unix())
		// if the response is not received in 10 seconds remove it from connected nodes
		// handle node failure
		if curr_time-send_time > 10 {
			if _, id := pending_connection_time[node]; id {
				delete(pending_connection_time, node)
				net_platform.AddToPreviousNodes(node)
				net_platform.RemoveNodeWithAddress(node)
			}
		}
	}
}

func Synchronize(net_platform *NetworkPlatform) {
	prev_time_send := uint64(time.Now().Unix())
	prev_time_check := uint64(time.Now().Unix())

	for _, node := range net_platform.Connected_nodes {
		SendEchoMessage(node.GetAddressString(), net_platform)
	}
	for true {
		curr_time := uint64(time.Now().Unix())

		// send every 20 sec
		if curr_time-prev_time_send > 20 {
			prev_time_send = curr_time
			for _, node := range net_platform.Connected_nodes {
				SendEchoMessage(node.GetAddressString(), net_platform)
			}
		}

		// check every 10 sec
		if curr_time-prev_time_check > 10 {
			prev_time_check = curr_time
			CheckForResponse(net_platform)
		}
	}
}
