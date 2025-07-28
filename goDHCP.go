package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	utils "github.com/jgib/utils"
)

var debug bool = false

type dhcpData struct {
	configFile string
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

type jsonData struct {
	ServerPort uint16 `json:"serverPort"`
	ClientPort uint16 `json:"clientPort"`
	PoolStart  string `json:"poolStart"`
	PoolEnd    string `json:"poolEnd"`
	OP         byte   `json:"OP"`
	HTYPE      byte   `json:"HTYPE"`
	HOPS       byte   `json:"HOPS"`
	XID        byte   `json:"XID"`
	FLAGS      byte   `json:"FLAGS"`
	YIADDR     string `json:"YIADDR"`
	SIADDR     string `json:"SIADDR"`
	GIADDR     string `json:"GIADDR"`
	CHADDR     string `json:"CHADDR"`
	SNAME      string `json:"SNAME"`
	FILE       string `json:"FILE"`
}

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
	data.configFile = "../config/goDHCP.json"
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
			fmt.Printf("  [-c | --config \"PATH\"]    Specify the JSON configuration file to use, Default: %s\n", data.configFile)
			fmt.Printf(" <-ps | --poolstart A.B.C.D>  First IP in the pool.\n")
			fmt.Printf(" <-pe | --poolend A.B.C.D>    Last IP in the pool.\n")
			fmt.Printf(" [-sp | --serverport UINT16]  Server port, Default %d.\n", data.serverPort)
			fmt.Printf(" [-cp | --clientport UINT16]  Client port, Default %d.\n", data.clientPort)

			os.Exit(0)
		}

		if (arg == "-c" || arg == "--config") && i+1 < len(args) {
			data.configFile = args[i+1]
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

		utils.Debug(arg, debug)
	}

	_, err = os.Stat(data.configFile)
	utils.Er(err)
	file, err := os.ReadFile(data.configFile)
	utils.Er(err)
	var jsonConfig jsonData
	json.NewDecoder(bytes.NewBuffer(file)).Decode(&jsonConfig)

	if data.poolStart == "" {
		data.poolStart = jsonConfig.PoolStart
	}
	if data.poolEnd == "" {
		data.poolEnd = jsonConfig.PoolEnd
	}

	startIp, err := utils.Ip2Uint32(data.poolStart)
	utils.Er(err)
	stopIp, err := utils.Ip2Uint32(data.poolEnd)
	utils.Er(err)
	if startIp > stopIp {
		utils.Er(fmt.Errorf("pool end IP [%s] cannot be larger than pool start IP [%s]", data.poolEnd, data.poolStart))
	}

	PrintData()

}
