package utils

import (
	"encoding/json"
	"os"

	"github.com/mmoehabb/luci/types"
)

func LoadDefaultConfig() types.Config {
	const configPath = "luci.config.json"

	// Open and read the configuration file
	_, err := os.Open(configPath)
	if err != nil {
		panic("luci.config.json file not found!")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	// Parse the json data and perform the action the user passes in the arguments
	c := types.Config{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}

	// Convert map[string]any maps with value keys to AnnotatedActions
	shellc := *GetShellConfig(c)
	digForAnnotatedActions(shellc)

	return c
}

func digForAnnotatedActions(m map[string]any) {
	for k, v := range m {
		switch v := v.(type) {
		case map[string]any:
			if v["value"] == nil {
				digForAnnotatedActions(v)
				continue
			}
			annAct := mapToAnnotatedAction(v)
			m[k] = annAct
		}
	}
}

func mapToAnnotatedAction(m map[string]any) types.AnnotatedAction {
	title := ""
	if m["title"] != nil {
		title = m["title"].(string)
	}

	description := ""
	if m["description"] != nil {
		description = m["description"].(string)
	}

	return types.AnnotatedAction{
		Title:       title,
		Description: description,
		Value:       m["value"],
	}
}
