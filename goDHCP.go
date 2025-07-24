package main

import (
	"fmt"

	utils "github.com/jgib/utils"
)

var debug bool = false

type dhcpData struct {
	poolStart  string
	poolEnd    string
	serverPort uint16
	clientPort uint16

	// OPTIONS
	op      byte
	htype   byte
	hlen    byte
	hops    byte
	xid     uint32
	secs    uint16
	flags   uint16
	ciaddr  uint32
	yiaddr  uint32
	siaddr  uint32
	giaddr  uint32
	chaddr  []byte // 16 Bytes
	sname   []byte // 64 Bytes
	file    []byte // 128 Bytes
	options []byte // Variable Bytes
}

var data dhcpData

func main() {
	data.serverPort = 67
	data.clientPort = 68
	data.op = 1
	data.htype = 1
	data.hlen = 6
	data.chaddr = make([]byte, 16)
	data.sname = make([]byte, 64)
	data.file = make([]byte, 128)
	data.options = []byte{}

	args, err := utils.GetArgs(0)
	utils.Er(err)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-h" || arg == "--help" {
			fmt.Printf("A pure Go DHCP server in accordance with RCC 2131 and RFC 2132.\n\n")

			fmt.Printf(" [-h | --help]  Print this help message.\n")

			fmt.Printf("\nServer Parameters:\n")
			fmt.Printf(" <-ps | --poolstart A.B.C.D>  First IP in the pool.\n")
			fmt.Printf(" <-pe | --poolend A.B.C.D>    Last IP in the pool.\n")
			fmt.Printf(" [-sp | --serverport UINT16]  Server port, Default %d.\n", data.serverPort)
			fmt.Printf(" [-cp | --clientport UINT16]  Client port, Default %d.\n", data.clientPort)

			fmt.Printf("\nDHCP Packet Overrides:\n")
			fmt.Printf(" [--op UINT8]       Message op code / message type, 1 = BOOTREQUEST, 2 = BOOTREPLY, Default: %d.\n", data.op)
			fmt.Printf(" [--htype UINT8]    Hardware address type, 1 = 10mb ethernet, see ARP section in \"Assigned Numbers RFC\", Default: %d.\n", data.htype)
			fmt.Printf(" [--hlen UINT8]     Hardware address length, 6 for 10mb ethernet, Default: %d.\n", data.hlen)
			fmt.Printf(" [--hops UINT8]     Client sets to 0, used by relay agents when booting via relay agent, Default: %d.\n", data.hops)
			fmt.Printf(" [--xid UINT32]     Transaction ID, random number chosen by the client, Default: %d.\n", data.xid)
			fmt.Printf(" [--flags]          Sets the broadcast bit, 0x%02X, Default: 0x%02X.\n", 0b1000000000000000, data.flags)
			fmt.Printf(" [--yiaddr A.B.C.D] 'your' (client) IP address, Default: %d.\n", data.yiaddr)
			fmt.Printf(" [--siaddr A.B.C.D] IP address of next server to use in bootstrap, Default: %d.\n", data.siaddr)
			fmt.Printf(" [--giaddr A.B.C.D] Relay agent IP address, used in booting via relay agent, Default: %d.\n", data.giaddr)
			fmt.Printf(" [--chaddr 16BYTES] Client hardware address, given in hexadecimal format, Default:%s\n", utils.WalkByteSlice(data.chaddr))
			fmt.Printf(" [--sname STRING]   Optional server host name, null terminated string, Max of 64 bytes.\n")
			fmt.Printf(" [--file STRING]    Boot file name, null terminated string, Max of 128 bytes.\n")

			fmt.Printf("\nDHCP Options:\n")
			fmt.Printf(" RFC 1497 Vendor Extensions:\n")
			fmt.Printf("  [-0]         Pad.\n")
			fmt.Printf("  [-255]       End.\n")
			fmt.Printf("  [-1 A.B.C.D] Subnet Mask.\n")
			fmt.Printf("  ...\n")
			fmt.Printf(" IP Layer Parameters per Host:\n")
			fmt.Printf(" IP Layer Parameters per Interface:\n")
			fmt.Printf(" Link Layer Parameters per Interface:\n")
			fmt.Printf(" TCP Parameters:\n")
			fmt.Printf(" Application and Service Parameters:\n")
			fmt.Printf(" DHCP Extensions:\n")
		}

		if arg == "-v" || arg == "--verbose" {
			debug = true
		}

		utils.Debug(arg, debug)
	}
}
