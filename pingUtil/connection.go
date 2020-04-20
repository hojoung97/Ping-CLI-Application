package pingUtil

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"os"
)

var (
	ipv4Types = []string {"udp4", "ip4:icmp"}
	ipv6Types = []string {"udp6", "ip6:ipv6-icmp"}
)

// initiate a connection
func OpenConn(ipAddr net.IP, root bool) (*icmp.PacketConn, string) {
	var rootFlag int

	if root {
		rootFlag = 1
	} else {
		rootFlag = 0
	}

	if len(ipAddr.To4()) == net.IPv4len {
		conn, err := icmp.ListenPacket(ipv4Types[rootFlag], "0.0.0.0")

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR listening to packets %v\n", err)
			return nil, "IPv4"
		}

		conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
		return conn, "IPv4"

	} else if len(ipAddr) == net.IPv6len {
		conn, err := icmp.ListenPacket(ipv6Types[rootFlag], "::")

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR listening to packets %v\n", err)
			return nil, "IPv6"
		}

		conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)
		return conn, "IPv6"

	}
	return nil, ""
}