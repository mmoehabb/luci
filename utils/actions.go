package utils

import (
	"github.com/mmoehabb/luci/types"
)

type ActionNode struct {
	Key         string
	Title       string
	Description string
	IsGroup     bool
	Action      any
}

func CollectActions(config map[string]any) []ActionNode {
	var nodes []ActionNode

	for key, val := range config {
		switch v := val.(type) {
		case types.AnnotatedAction:
			isGroup := false
			if nested, ok := v.Value.(map[string]any); ok {
				if _, hasValue := nested["value"]; hasValue {
					isGroup = false
				} else {
					isGroup = true
				}
			}

			if !isGroup {
				node := ActionNode{
					Key:         key,
					Title:       v.Title,
					Description: v.Description,
					IsGroup:     false,
					Action:      v,
				}
				if node.Title == "" {
					node.Title = key
				}
				nodes = append(nodes, node)
			} else {
				node := ActionNode{
					Key:         key,
					Title:       v.Title,
					Description: v.Description,
					IsGroup:     true,
					Action:      nil,
				}
				if node.Title == "" {
					node.Title = key
				}
				nodes = append(nodes, node)
			}

		case map[string]any:
			if v["value"] != nil {
				switch v["value"].(type) {
				case map[string]any:
					node := ActionNode{
						Key:         key,
						Title:       MapToAnnotatedAction(v).Title,
						Description: MapToAnnotatedAction(v).Description,
						IsGroup:     true,
						Action:      nil,
					}
					if node.Title == "" {
						node.Title = key
					}
					nodes = append(nodes, node)
				case []string, string:
					ann := MapToAnnotatedAction(v)
					node := ActionNode{
						Key:         key,
						Title:       ann.Title,
						Description: ann.Description,
						IsGroup:     false,
						Action:      ann,
					}
					if node.Title == "" {
						node.Title = key
					}
					nodes = append(nodes, node)
				}
			} else {
				node := ActionNode{
					Key:         key,
					Title:       key,
					Description: "",
					IsGroup:     true,
					Action:      nil,
				}
				nodes = append(nodes, node)
			}

		case string:
			node := ActionNode{
				Key:         key,
				Title:       key,
				Description: "",
				IsGroup:     false,
				Action:      v,
			}
			nodes = append(nodes, node)

		case []string:
			node := ActionNode{
				Key:         key,
				Title:       key,
				Description: "",
				IsGroup:     false,
				Action:      v,
			}
			nodes = append(nodes, node)
		}
	}

	return nodes
}

const BackKey = "__back__"
