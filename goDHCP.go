package main

import (
	utils "github.com/jgib/utils"
)

var debug bool = false

func main() {
	args, err := utils.GetArgs(0)
	utils.Er(err)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-v" || arg == "--verbose" {
			debug = true
		}

		utils.Debug(arg, debug)
	}
}
