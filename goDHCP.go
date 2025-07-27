package main

import (
	"fmt"
	"os"
	"strconv"

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

func PrintData() {
	utils.Debug(fmt.Sprintf("POOL START: %s", data.poolStart), debug)
	utils.Debug(fmt.Sprintf("POOL STOP:  %s", data.poolEnd), debug)
	utils.Debug(fmt.Sprintf("SERVER PORT:%d", data.serverPort), debug)
	utils.Debug(fmt.Sprintf("CLIENT PORT:%d", data.clientPort), debug)
	utils.Debug(fmt.Sprintf("OP:         %d", data.op), debug)
	utils.Debug(fmt.Sprintf("HTYPE:      %d", data.htype), debug)
	utils.Debug(fmt.Sprintf("HLEN:       %d", data.hlen), debug)
	utils.Debug(fmt.Sprintf("HOPS:       %d", data.hops), debug)
	utils.Debug(fmt.Sprintf("XID:        %d", data.xid), debug)
	utils.Debug(fmt.Sprintf("SECS:       %d", data.secs), debug)
	utils.Debug(fmt.Sprintf("FLAGS:      %d", data.flags), debug)
	utils.Debug(fmt.Sprintf("CIADDR:     %d", data.ciaddr), debug)
	utils.Debug(fmt.Sprintf("YIADDR:     %d", data.yiaddr), debug)
	utils.Debug(fmt.Sprintf("SIADDR:     %d", data.siaddr), debug)
	utils.Debug(fmt.Sprintf("GIADDR:     %d", data.giaddr), debug)
	utils.Debug(fmt.Sprintf("CHADDR:\n%s", utils.WalkByteSlice(data.chaddr)), debug)
	utils.Debug(fmt.Sprintf("SNAME:      %s", data.sname), debug)
	utils.Debug(fmt.Sprintf("FILE:       %s", data.file), debug)
	utils.Debug(fmt.Sprintf("OPTIONS:\n%s", utils.WalkByteSlice(data.options)), debug)
}

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

			fmt.Printf("\nDHCP Options:")
			fmt.Printf("\n RFC 1497 Vendor Extensions:\n")
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

			fmt.Printf("\n IP Layer Parameters per Host:\n")
			fmt.Printf("  [-19]                 IP Forwarding: Specifies whether the client should configure its IP layer for packet forwarding.\n")
			fmt.Printf("  [-20]                 Non-Local Source Routing: Specifies whether the client should configure its IP layer to allow forwarding of\n")
			fmt.Printf("                        datagrams with non-local source routes.\n")
			fmt.Printf("  [-21 A.B.C.D M.M.M.M] Policy Filter: Specifies policy filters  for non-local source routing. The filters consist of a list of IP\n")
			fmt.Printf("                        addresses and masks which specify destination/mask pairs with which to filter incoming source routes.")
			fmt.Printf("  ...\n")
			fmt.Printf("\n IP Layer Parameters per Interface:\n")
			fmt.Printf("\n Link Layer Parameters per Interface:\n")
			fmt.Printf("\n TCP Parameters:\n")
			fmt.Printf("\n Application and Service Parameters:\n")
			fmt.Printf("\n DHCP Extensions:\n")

			os.Exit(0)
		}

		if arg == "-v" || arg == "--verbose" {
			debug = true
		}

		if (arg == "-ps" || arg == "--poolstart") && i+1 < len(args) {
			_, err := utils.Ip2Uint32(args[i+1])
			utils.Er(err)
			data.poolStart = args[i+1]
		}

		if (arg == "-pe" || arg == "--poolend") && i+1 < len(args) {
			_, err := utils.Ip2Uint32(args[i+1])
			utils.Er(err)
			data.poolEnd = args[i+1]
		}

		if (arg == "-sp" || arg == "--serverport") && i+1 < len(args) {
			tmp, err := utils.Port2Uint16(args[i+1])
			utils.Er(err)
			data.serverPort = tmp
		}

		if (arg == "-cp" || arg == "--clientport") && i+1 < len(args) {
			tmp, err := utils.Port2Uint16(args[i+1])
			utils.Er(err)
			data.clientPort = tmp
		}

		if arg == "--op" && i+1 < len(args) {
			tmp, err := strconv.ParseUint(args[i+1], 10, 8)
			utils.Er(err)
			data.op = byte(tmp)
		}

		if arg == "--htype" && i+1 < len(args) {
			tmp, err := strconv.ParseUint(args[i+1], 10, 8)
			utils.Er(err)
			data.htype = byte(tmp)
		}

		if arg == "--hlen" && i+1 < len(args) {
			tmp, err := strconv.ParseUint(args[i+1], 10, 8)
			utils.Er(err)
			data.hlen = byte(tmp)
		}

		if arg == "--hops" && i+1 < len(args) {
			tmp, err := strconv.ParseUint(args[i+1], 10, 8)
			utils.Er(err)
			data.hops = byte(tmp)
		}

		if arg == "--xid" && i+1 < len(args) {
			tmp, err := strconv.ParseUint(args[i+1], 10, 32)
			utils.Er(err)
			data.xid = uint32(tmp)
		}

		if arg == "--flags" {
			data.flags = 0b1000000000000000
		}

		if arg == "--yiaddr" && i+1 < len(args) {
			tmp, err := utils.Ip2Uint32(args[i+1])
			utils.Er(err)
			data.yiaddr = tmp
		}

		if arg == "--siaddr" && i+1 < len(args) {
			tmp, err := utils.Ip2Uint32(args[i+1])
			utils.Er(err)
			data.siaddr = tmp
		}

		if arg == "--giaddr" && i+1 < len(args) {
			tmp, err := utils.Ip2Uint32(args[i+1])
			utils.Er(err)
			data.giaddr = tmp
		}

		if arg == "--chaddr" && i+1 < len(args) {
			for j, k := 0, 0; j < len(args[i+1]); j, k = j+2, k+1 {
				if j+1 < len(args[i+1]) {
					tmp, err := strconv.ParseUint(args[i+1][j:j+2], 16, 8)
					utils.Er(err)
					data.chaddr[k] = byte(tmp)
				}
			}
		}

		if arg == "--sname" && i+1 < len(args) {
			for j := 0; j < len(args[i+1]); j++ {
				if j >= len(data.sname) {
					break
				}

				data.sname[j] = args[i+1][j]
			}
			data.sname[len(data.sname)-1] = 0
		}

		if arg == "--file" && i+1 < len(args) {
			for j := 0; j < len(args[i+1]); j++ {
				if j >= len(data.file) {
					break
				}

				data.file[j] = args[i+1][j]
			}
			data.file[len(data.file)-1] = 0
		}

		if arg == "-0" {
			data.options = append(data.options, 0)
		}

		if arg == "-255" {
			data.options = append(data.options, 255)
		}

		if (arg == "-1" || arg == "-2" || arg == "-3" || arg == "-4" || arg == "-5" || arg == "-6" || arg == "-7" || arg == "-8" || arg == "-9" || arg == "-10" ||
			arg == "-11") && i+1 < len(args) {
			tmp, err := strconv.ParseUint(args[i+1], 10, 32)
			utils.Er(err)
			data.options = append(data.options, byte(tmp>>24))
			data.options = append(data.options, byte(tmp>>16))
			data.options = append(data.options, byte(tmp>>8))
			data.options = append(data.options, byte(tmp))
		}

		utils.Debug(arg, debug)
	}

	PrintData()

	startIp, err := utils.Ip2Uint32(data.poolStart)
	utils.Er(err)
	stopIp, err := utils.Ip2Uint32(data.poolEnd)
	utils.Er(err)
	if startIp > stopIp {
		utils.Er(fmt.Errorf("pool end IP [%s] cannot be larger than pool start IP [%s]", data.poolEnd, data.poolStart))
	}
}
