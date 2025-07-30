package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	utils "github.com/jgib/utils"
)

/* Example Go Route:

func testFunc(j string, wg *sync.WaitGroup) {
	fmt.Printf("TEST:%s\r\n", j)
	defer wg.Done()
}

	var wg sync.WaitGroup
	wg.Add(len(ips))
	for i := 0; i < len(n); i++ {


		go testFunc(ip, &wg)

	}
	wg.Wait()
*/

var debug bool = false

type dhcpData struct {
	configFile string
	poolStart  string `json:"poolStart"`
	poolEnd    string `json:"poolEnd"`
	serverPort uint16 `json:"serverPort"`
	clientPort uint16 `json:"clientPort"`

	// OPTIONS
	op          byte   `json:"OP"`
	htype       byte   `json:"HTYPE"`
	hlen        byte   `json:"HLEN"`
	hops        byte   `json:"HOPS"`
	xid         uint32 `json:"XID"`
	secs        uint16 `json:"SECS"`
	flags       uint16 `json:"FLAGS"`
	ciaddr      string `json:"CIADDR"`
	yiaddr      string `json:"YIADDR"`
	siaddr      string `json:"SIADDR"`
	giaddr      string `json:"GIADDR"`
	chaddr      string `json:"CHADDR"`
	chaddrBytes []byte // 16 Bytes
	sname       string `json:"SNAME"` // 64 Bytes
	file        string `json:"FILE"`  // 128 Bytes
	options     []byte // Variable Bytes
	opt1        string `json:"OPT1"`
	opt2        uint32 `json:"OPT2"`
	opt3        string `json:"OPT3"`
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
	utils.Debug(fmt.Sprintf("CHADDR:\n%s", data.chaddr), debug)
	utils.Debug(fmt.Sprintf("SNAME:      %s\n", data.sname), debug)
	utils.Debug(fmt.Sprintf("FILE:       %s\n", data.file), debug)
	utils.Debug(fmt.Sprintf("OPTIONS:\n%s", utils.WalkByteSlice(data.options)), debug)

}

func main() {
	filePath, err := os.Executable()
	utils.Er(err)
	data.configFile = fmt.Sprintf("%s/../config/goDHCP.json", filepath.Dir(filePath))
	data.serverPort = 67
	data.clientPort = 68
	data.op = 1
	data.htype = 1
	data.hlen = 6
	data.chaddrBytes = make([]byte, 16)
	data.options = []byte{}

	args, err := utils.GetArgs(0)
	utils.Er(err)

	_, err = os.Stat(data.configFile)
	utils.Er(err)
	file, err := os.ReadFile(data.configFile)
	utils.Er(err)
	json.NewDecoder(bytes.NewBuffer(file)).Decode(&data)

	data.chaddr = []byte(jsonConfig.CHADDR) // convert from hex
	for i := 0; i < len([]byte(jsonConfig.SNAME)); i++ {
		if i >= len(data.sname) {
			break
		}
		data.sname[i] = []byte(jsonConfig.SNAME)[i]
	}
	for i := 0; i < len([]byte(jsonConfig.FILE)); i++ {
		if i >= len(data.file) {
			break
		}
		data.file[i] = []byte(jsonConfig.FILE)[i]
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-h" || arg == "--help" {
			fmt.Printf("A pure Go DHCP server in accordance with RCC 2131 and RFC 2132.\n\n")

			fmt.Printf(" [-h | --help]  Print this help message.\n")

			fmt.Printf("\nServer Parameters:\n")
			fmt.Printf(" [ -c | --config \"PATH\" ]     Specify the JSON configuration file to use.\n")
			fmt.Printf("                                Default: %s\n", data.configFile)
			fmt.Printf(" < -ps | --poolstart A.B.C.D >  First IP in the pool.\n")
			fmt.Printf(" < -pe | --poolend A.B.C.D >    Last IP in the pool.\n")
			fmt.Printf(" [ -sp | --serverport UINT16 ]  Server port, Default %d.\n", data.serverPort)
			fmt.Printf(" [ -cp | --clientport UINT16 ]  Client port, Default %d.\n", data.clientPort)
			fmt.Printf("\n")
			fmt.Printf("Specific DHCP options, as well as the options listed above, are configured in\n")
			fmt.Printf("the file 'config/example.json'\n")
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
	}

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

	data.sname[len(data.sname)-1] = 0
	data.file[len(data.file)-1] = 0

	data.chaddrBytes, err = hex.DecodeString(data.chaddr)
	utils.Er(err)

	utils.Debug(utils.WalkByteSlice(data.chaddrBytes), debug)
	fmt.Printf("CHADDR [%s]\n", utils.WalkByteSlice(data.chaddrBytes))

	PrintData()
}
