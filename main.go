package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mmoehabb/luci/utils"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("LUCI version %s\n", "0.0.1")
		os.Exit(0)
	}

	c := utils.LoadDefaultConfig()
	if len(os.Args) < 2 {
		utils.PrintUsage(c)
	} else {
		utils.Act(c, os.Args[1:])
	}
}
