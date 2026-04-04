package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mmoehabb/luci/utils"
)

func main() {
	listFlag := flag.Bool("list", false, "List actions without interactive selection")
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("LUCI version %s\n", "0.0.6")
		os.Exit(0)
	}

	c := utils.LoadDefaultConfig()
	if *listFlag {
		utils.PrintUsage(c)
	} else if len(os.Args) < 2 {
		utils.PrintInteractiveUsage(c)
	} else {
		utils.Act(c, os.Args[1:])
	}
}
