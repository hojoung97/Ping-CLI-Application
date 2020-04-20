package pingUtil

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"os"
)

func CreateEchoRequest (ipType string, seq int) []byte{
	var msgType icmp.Type

	// check if IPv4 or IPv6
	if ipType == "IPv4" {
		msgType = ipv4.ICMPTypeEcho
	} else {
		msgType = ipv6.ICMPTypeEchoRequest
	}

	// initialize ICMP message
	msg := icmp.Message{
		Type: msgType,
		Code: 0,	// code 0 for echo request
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff,
			Seq: seq,
			Data: []byte("Ping: Echo Request send"),
		},
	}

	// encode into bytes
	encodedMsg, err := msg.Marshal(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR while encoding ICMP message: %v\n", err)
	}

	return encodedMsg
}
