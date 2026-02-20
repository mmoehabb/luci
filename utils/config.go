package utils

import (
	"bufio"
	"io/fs"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/mmoehabb/luci/types"
)

const InitConfigStr = `
title = "Hello World!"
description = "Luci config hello world example."

[bash.run]
example = "echo Hello World!"

[zshell.run]
example = "echo Hello World!"

[bat.run]
example = "echo Hello World!"
`

// LoadDefaultConfig loads the application configuration from luci.config.toml.
// If the configuration file does not exist, it creates a default configuration
// file with example settings. The function returns a populated Config struct
// that contains all shell-specific settings and metadata.
func LoadDefaultConfig() types.Config {
	const configPath = "luci.config.toml"

	// Open and read the configuration file
	_, err := os.Open(configPath)
	if err != nil {
		color.Red("luci.config.toml file not found!")
		println()
		color.Yellow("Create default config? [y/N] ")
		if readApproval() {
			createDefaultConfig()
		} else {
			log.Fatalln("You have to write luci.config.toml file in order to use luci.")
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	// Parse the json data and perform the action the user passes in the arguments
	c := types.Config{}
	Must(toml.Unmarshal(data, &c))
	return c
}

func createDefaultConfig() {
	const configPath = "luci.config.toml"
	err := os.WriteFile(configPath, []byte(InitConfigStr), fs.ModePerm)
	if err != nil {
		panic("luci.config.toml creation failed!")
	}
	log.Println("✓ luci.config.toml has been created")
}

func readApproval() bool {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		panic(err)
	}
	return char == 'Y' || char == 'y'
}
