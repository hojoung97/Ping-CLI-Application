package main

import (
	"Cloudflare2020/pingUtil"
	"flag"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

// these are defined in an internal package so I manually defined here
// https://godoc.org/golang.org/x/net/internal/iana
var (
	ProtocolICMP = 1
	ProtocolIPv6ICMP = 58
)

func main () {

	// flag options
	timeout := flag.Duration("t", time.Second*10, "maximum wait time before exiting for no response")
	interval := flag.Duration("i", time.Second, "interval between each echo request")
	count := flag.Int("c", -1, "specified amount of echo requests before exiting")
	packetSize := flag.Int("s", 56, "size of each packet in bytes")
	root := flag.Bool("root", false, "whether given root privilege or not" +
														"\nIf using with 'sudo' MUST set this flag")
	help := flag.Bool("h", false, "print helpful usage statements")
	flag.Parse()

	// print usage message regarding flags
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// hostname or IP address input is a REQUIRED parameter, abort if not given
	if flag.NArg() != 1 {
		fmt.Printf(
			"Usage: %s [-t timeout] [-i interval] [-c count] [-s packetSize] [-root root]" +
				" [-h help] <hostname or IP address>\n", os.Args[0])
		os.Exit(1)
	}

	// Prepare target IP address according to the given input format
	addr, err := net.ResolveIPAddr("ip", flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR in ResolveIPAddr for given input target: %v\n", err)
	}

	// make sure destination is in UDP connection if root privilege is not given
	// but if root privilege is given just use IP connection
	var dst net.Addr
	if !*root {
		dst = &net.UDPAddr{
			IP:   addr.IP,
			Port: 0,
			Zone: addr.Zone,
		}
	} else {
		dst = addr
	}

	// starting message
	fmt.Printf("PING %s (%s): %d data bytes\n", flag.Arg(0), addr, *packetSize)

	// open connection
	conn, ipType := pingUtil.OpenConn(addr.IP, *root)
	if conn == nil {
		fmt.Fprint(os.Stderr, "ERROR initiating connection\n")
		os.Exit(1)
	}
	defer conn.Close()		// close the connection when finished with everything

	var stat pingUtil.Statistic
	stat.Dst = flag.Arg(0)

	// listen to ctrl c keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			stat.PrintStats()
			os.Exit(0)
		}
	}()

	// infinite loop that sends echo request to target
	for i := 0;;i++ {
		// if count flag is specified only run until specified amount of echo request/reply is sent/received
		if  i == *count {
			break
		}

		// create ICMP echo request message
		msg := pingUtil.CreateEchoRequest(ipType, i)

		// send
		start := time.Now()
		numBytes, err := conn.WriteTo(msg, dst)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR sending echo request to destination target1: %v\n", err)
			os.Exit(1)
		} else if numBytes != len(msg) {
			fmt.Fprintf(os.Stderr, "ERROR sending echo request to destintation target2: %v\n", err)
			os.Exit(1)
		}

		// update statistics
		stat.PackTrans += 1

		// wait to receive a reply
		reply := make([]byte, *packetSize+8)
		err = conn.SetReadDeadline(time.Now().Add(*timeout))
		if err != nil {
			fmt.Printf("ERROR while setting timeout for echo request %v\n", err)
		}

		var peerIPv4 *ipv4.ControlMessage
		var peerIPv6 *ipv6.ControlMessage

		var protocol int	// Protocol Type

		// read the received reply
		if ipType == "IPv4" {
			numBytes, peerIPv4, _, err = conn.IPv4PacketConn().ReadFrom(reply)
			protocol = ProtocolICMP
		} else if ipType == "IPv6" {
			numBytes, peerIPv6, _, err = conn.IPv6PacketConn().ReadFrom(reply)
			protocol = ProtocolIPv6ICMP
		}
		if err != nil {
			if strings.HasSuffix(err.Error(), "timeout") {
				fmt.Printf("Request timeout for icmp_seq %d\n", i)
				continue
			} else {
				fmt.Printf("ERROR while receiving echo reply: %v\n", err)
				continue
			}
		}

		// update statistics
		stat.PackRecv += 1

		// round-trip time
		rtt := time.Since(start)
		stat.Rtts = append(stat.Rtts, rtt.Seconds())

		// parse the message replied from the target
		echoReply, err := icmp.ParseMessage(protocol, reply[:numBytes])
		if err != nil {
			fmt.Printf("ERROR while parsing echo reply: %v\n", err)
		}

		// print statements according to the IP version
		switch echoReply.Type {
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("%v bytes from %v: icmp_seq=%d ttl=%v time=%v\n",
				numBytes, peerIPv4.Src, i, peerIPv4.TTL, rtt)
		case ipv6.ICMPTypeEchoReply:
			fmt.Printf("%v bytes from %v: icmp_seq=%d ttl=%v time=%v\n",
				numBytes, peerIPv6.Src, i, peerIPv6.HopLimit, rtt)
		default:
			fmt.Printf("Unexpeted reply received")
		}

		// give some interval before next consecutive echo requests
		time.Sleep(*interval)
	}
}