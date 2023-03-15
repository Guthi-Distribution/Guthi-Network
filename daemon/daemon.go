package daemon

/*
#cgo LDFLAGS: -L./ -ldaemon
#include "daemon_api.h"
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"
)

// Replicated for now
const (
	GetFile              = 0
	CheckIfInCache       = 1
	RequestFileMetadata  = 2
	NoSuchResourceExists = 3
	TrackedFileChanged   = 4 // Response : followed by (modification time) File Name
	TrackThisFile        = 5 // Request  : same as above
	EchoMessage          = 6
	Continuation         = 7 // Response : Continuation of previous message
	TrackThisFolder      = 8
)

type DaemonHandle struct {
	handle        C.Handle
	FormatChannel chan MsgFormat
}

type MsgFormat struct {
	Msg_type    uint8
	Msg_len     uint32
	Msg_content []byte
}

// Bruh, OOP like -> forced shit
func BuildMessage(mtype uint8, mlen uint32, mcontent []byte) MsgFormat {
	return MsgFormat{Msg_type: mtype, Msg_len: mlen, Msg_content: mcontent}
}

func InitializeDaemon() *DaemonHandle {
	handle := DaemonHandle{handle: C.GetDaemonHandle()}
	// vals := [13]byte{0x48, 0x55, 0x54, 0x49, 0x4A, 0x06, 0x05, 0x00, 'H', 'e', 'l', 'l', 'o'}
	// C.PostMessageToDaemon(handle.handle, (*C.uchar)(unsafe.Pointer(&vals)), 13)
	// Wait for messages
	// val := new([512]byte)

	// len := uint32(0)
	// C.GetMessageFromDaemon(handle.handle, (*C.uchar)(unsafe.Pointer(&val[0])), (*C.uint)(&len), 512)
	// fmt.Println("Message received from daemon : ", len, string(val[:len]))
	handle.FormatChannel = make(chan MsgFormat)
	// Check connection for echo back reply

	// Listen continuously for changes
	// for i := 0; i < 10; i++ {
	// 	C.PostMessageToDaemon(handle.handle, (*C.uchar)(unsafe.Pointer(&vals)), 13)
	// 	// Wait for back reply
	// 	// TODO :: Implement the non blocking read call, write can block all they want but blocking on read is not acceptable
	// 	C.GetMessageFromDaemon(handle.handle, (*C.uchar)(unsafe.Pointer(&val[0])), (*C.uint)(&len), 512)
	// 	fmt.Println("Message received from daemon : ", len, string(val[:len]))
	// }
	return &handle
}

func PrepareMessageForDaemon(msg []byte, msg_type int16 /* No Enum LOL*/, data []byte, length int16) uint32 {
	// Magic bytes
	msg[0] = 0x48
	msg[1] = 0x55
	msg[2] = 0x54
	msg[3] = 0x49
	msg[4] = 0x4A

	// type of message and length
	msg[5] = uint8(msg_type)
	msg[6] = uint8(length & 0xFF)
	msg[7] = uint8(length >> 8)

	// append bytes, but how?
	copy(msg[8:], data[:length])
	return uint32(length + 8)

	// This completes the message format
}

func VerifyMessageFromDaemon(msg []byte, length uint) bool {
	if length < 5 {
		return false
	}

	is_equal := (msg[0] == 0x48) &&
		(msg[1] == 0x55) &&
		(msg[2] == 0x54) &&
		(msg[3] == 0x49) &&
		(msg[4] == 0x4A)

	fmt.Println("Is equal : ", is_equal)
	return is_equal
}

func ParseMessage(msg []byte, length uint) MsgFormat {
	// First the verficiation phase
	if !VerifyMessageFromDaemon(msg, length) {
		fmt.Println("MessageParser : Magic bytes verification failed")
		os.Exit(-3)
	}
	rem_msg := msg[5:]
	rem_len := length - 5

	if rem_len < 3 {
		fmt.Println("Incomplete message")
		os.Exit(-4)
	}
	msg_len := uint32(rem_msg[2])<<8 + uint32(rem_msg[1])

	return MsgFormat{Msg_type: rem_msg[0], Msg_len: uint32(msg_len), Msg_content: rem_msg[3:]}
}

// TODO :: Implement streaming
func PollMessagesFromDaemon(daemon DaemonHandle) {
	msg := new([512]byte) // TODO :: Reduce this allocation

	// Do not block, I hate blocking
	read_bytes := C.ReadNonBlockingMessageFromDaemon(daemon.handle, (*C.uchar)(unsafe.Pointer(&msg[0])), 512)
	if read_bytes == -1 {
		fmt.Println("Connection closed/terminated from the daemon")
		os.Exit(-1)
	}

	// fmt.Println("Read : ", read_bytes)
	if read_bytes == 0 {
		// If it is 0 continue
		// For research purpose sleep here
		time.Sleep(time.Millisecond * 100)
		// fmt.Println("Non blocking on progress")
	} else {
		real_msg := ParseMessage(msg[:read_bytes], uint(read_bytes))
		log.Println(real_msg.Msg_type)
		daemon.FormatChannel <- real_msg
	}
}

func SendFormattedMessage(dae DaemonHandle, format MsgFormat) bool {
	// Allocate enough memory first
	send_buf := make([]byte, format.Msg_len+8)
	send_length := PrepareMessageForDaemon(send_buf, int16(format.Msg_type), format.Msg_content, int16(format.Msg_len))
	C.PostMessageToDaemon(dae.handle, (*C.uchar)(unsafe.Pointer(&send_buf[0])), C.uint(send_length))
	return true
}
