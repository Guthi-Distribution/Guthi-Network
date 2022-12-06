package ntp

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"time"
)

const (
	NTP_PORT            = 123
	VERSION             = 4
	TIMEOUT_IN_SECONDS  = 1
	NTP_TIMESTAMP_DELTA = 2208988800
)

var (
	TimeServers     []string
	TimeServerIndex int  = 0
	ValuesInit      bool = false
)

type ntpTimestamp struct {
	Seconds  uint32
	Fraction uint32
}

type ntpShort struct {
	Seconds  uint16
	Fraction uint16
}

var GlobalClock time.Time

type NtpPacket struct {
	LiVnMd            uint8        `json:"li_vn_md"`
	Stratum           uint8        `json:"stratum"`
	Poll              uint8        `json:"poll"`
	Precision         uint8        `json:"precision"`
	RootDelay         ntpShort     `json:"root_delay"`
	RootDispersion    ntpShort     `json:"root_dispersion"`
	RefID             uint32       `json:"ref_id"`
	RefTimestamp      ntpTimestamp `json:"ref_timestamp"`
	OriginTimestamp   ntpTimestamp `json:"origin_timestamp"`
	RecvTimestamp     ntpTimestamp `json:"recv_timestamp"`
	TransmitTimestamp ntpTimestamp `json:"transmit_timestamp"`
}

// run this function in a goroutine, will update GlobalClock for use
// TODO: might need to use mutex? does go have something like that?
func SyncTime() {
	for {
		var utcTime *time.Time
		// TODO: a slight delay in switching between faulty servers?
		for time, err := FetchUTCTime(); ; SwitchTimeServer() {
			if err == nil {
				utcTime = time
				break
			}
		}
		ticks := 1
		for {
			time.Sleep(time.Nanosecond)
			GlobalClock = time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(), utcTime.Hour(), utcTime.Minute(), utcTime.Second()+ticks, utcTime.Nanosecond(), utcTime.Location())
			ticks++
			// fmt.Println("global clock", GlobalClock.UTC()) // if you need to log time to something
			if ticks >= 10 {
				break
			}
		}
	}
}

func NTPToUnix(timestamp ntpTimestamp) uint32 {
	return timestamp.Seconds - NTP_TIMESTAMP_DELTA
}

func SetPacketParams(np *NtpPacket) {
	np.LiVnMd = 0b00010011 // leap invar utcTime *time.Timedicator = 0, version = 4, mode = 3
}

func ConnectToNTPServer() (net.Conn, error) {
	conn, err := net.Dial("udp", TimeServers[TimeServerIndex]+":123")
	if err != nil {
		return nil, err
	}
	err = conn.SetDeadline(time.Now().Add(TIMEOUT_IN_SECONDS * time.Second))
	return conn, err
}

func FetchUTCTime() (*time.Time, error) {
	// close connection after fetching utc time
	conn, err := ConnectToNTPServer()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	tp := NtpPacket{}
	SetPacketParams(&tp)
	err = binary.Write(conn, binary.BigEndian, tp)

	if err != nil {
		return nil, err
	}

	response := NtpPacket{}
	err = binary.Read(conn, binary.BigEndian, &response)
	if err != nil {
		return nil, err
	}

	// incorporate the returned nanoseconds value as well
	utcTime := time.Unix(int64(NTPToUnix(response.RecvTimestamp)), 0)
	return &utcTime, err
}

func LoadTimeServerInformation(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var payload map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}

	for _, v := range payload {
		TimeServers = append(TimeServers, v)
	}

	return nil
}

func SwitchTimeServer() {
	TimeServerIndex += 1
	if TimeServerIndex == len(TimeServers) {
		TimeServerIndex = 0
	}
}
