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
			fmt.Printf(" [--chaddr 16BYTES] Client hardware address, given in hexadecimal format,\n")
			fmt.Printf("                    Default:%s\n", utils.WalkByteSlice(data.chaddr))
			fmt.Printf(" [--sname STRING]   Optional server host name, null terminated string, Max of 64 bytes.\n")
			fmt.Printf(" [--file STRING]    Boot file name, null terminated string, Max of 128 bytes.\n")

			fmt.Printf("\nDHCP Options:\n")
			fmt.Printf(" RFC 1497 Vendor Extensions:\n")
			fmt.Printf("  [-0]          Pad.\n")
			fmt.Printf("  [-255]        End.\n")
			fmt.Printf("  [-1 A.B.C.D]  Subnet Mask: Specifies the client's subnet mask as per RFC 950.\n")
			fmt.Printf("  [-2 INT32]    Time Offset: Specifies the offset of the client's subnet in seconds from UTC. Expressed as 2's complement of INT32\n")
			fmt.Printf("  [-3 A.B.C.D]  Router: Specifies IP address for router on the client's subnet. Use parameter multiple times for multiple routers,\n")
			fmt.Printf("                should be in order of preference.\n")
			fmt.Printf("  [-4 A.B.C.D]  Time Server: Specifies time server available to the client. Use parameter multiple times for multiple servers,\n")
			fmt.Printf("                should be in order of preference.\n")
			fmt.Printf("  [-5 A.B.C.D]  Name Server: Specifies IEN 116 name servers avialable to the client. Use parameter multiple times for multiple\n")
			fmt.Printf("                servers, should be in order of preference.\n")
			fmt.Printf("  [-6 A.B.C.D]  Domain Name Server: Specifies DNS servers available to the client. Use parameter multiple times for multiple\n")
			fmt.Printf("                servers, should be in order of preference.\n")
			fmt.Printf("  [-7 A.B.C.D]  Log Server: Specifies MIT-LCS UDP log servers available to the client. Use parameter multiple times for multiple\n")
			fmt.Printf("                servers, should be in order of preference.\n")
			fmt.Printf("  [-8 A.B.C.D]  Cookie Server: Specifies RFC 865 cookie servers available to the client. Use parameter multiple times for multiple\n")
			fmt.Printf("                servers, should be in order of preference.\n")
			fmt.Printf("  [-9 A.B.C.D]  LPR Server: Specifies RFC 1179 line printer servers available to the client. Use parameter multiple times for multiple\n")
			fmt.Printf("                servers, should be in order of preference.\n")
			fmt.Printf("  [-10 A.B.C.D] Impress Server: Specifies Imagen Impress servers available to the client. Use parameter multiple times for multiple\n")
			fmt.Printf("                servers, should be in order of preference.\n")
			fmt.Printf("  [-11 A.B.C.D] Resource Location Server: Specifies Resource Location servers available to the client. Use parameter multiple times\n")
			fmt.Printf("                for multipleservers, should be in order of preference.\n")
			fmt.Printf("  [-12 STRING]  Host Name: Specifies the name of the client. \n")
			fmt.Printf("  [-13 UINT16]  Boot File Size: Specifies the length in 512-octet blocks of the default boot image for the client.\n")
			fmt.Printf("  [-14 STRING]  Merit Dump File: Specifies the path-name of a file to which the client's core image should be dumped in the event the\n")
			fmt.Printf("                client crashes.\n")
			fmt.Printf("  [-15 STRING]  Domain Name: Specifies the domain name that client should use when resolving hostnames via DNS.\n")
			fmt.Printf("  [-16 A.B.C.D] Swap Server: Specifies the IP address of the client's swap server.\n")
			fmt.Printf("  [-17 STRING]  Root Path: Specifies the path-name that contains the client's root disk.\n")
			fmt.Printf("  [-18 STRING]  Extensions Path: Specifies a file, retrievable via TFTP.\n")

			fmt.Printf(" IP Layer Parameters per Host:\n")
			fmt.Printf("  ...\n")
			fmt.Printf(" IP Layer Parameters per Interface:\n")
			fmt.Printf(" Link Layer Parameters per Interface:\n")
			fmt.Printf(" TCP Parameters:\n")
			fmt.Printf(" Application and Service Parameters:\n")
			fmt.Printf(" DHCP Extensions:\n")
		}

		if arg == "-v" || arg == "--verbose" {
			debug = true
		}

		if (arg == "-ps" || arg == "--poolstart") && i+1 < len(args) {
			data.poolstart = args[i+1] // validate ip
		}

		utils.Debug(arg, debug)
	}
}
