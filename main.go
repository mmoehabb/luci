package main

import (
	"os"

	"github.com/mmoehabb/luci/utils"
)

func main() {
	c := utils.LoadDefaultConfig()
	if len(os.Args) < 2 {
		utils.PrintUsage(c)
	} else {
		utils.Act(c, os.Args[1:])
	}
}
