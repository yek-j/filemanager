package plugins

import "fmt"

func GetPlugin(name string) (Plugin, error) {
	switch name {
	case "underscore_number":
		return &UnderscoreNumber{}, nil
	default:
		return nil, fmt.Errorf("unknown plugin: %s", name)
	}
}
