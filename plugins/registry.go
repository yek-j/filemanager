package plugins

import (
	"fmt"

	"github.com/yek-j/filemanager/config"
)

func GetPlugin(pluginCfg *config.PluginConfig) (Plugin, error) {
	switch pluginCfg.Name {
	case "underscore_number":
		return &UnderscoreNumber{pluginCfg: pluginCfg}, nil
	case "file_relocator":
		return &FileRelocator{pluginCfg: pluginCfg}, nil
	default:
		return nil, fmt.Errorf("unknown plugin: %s", pluginCfg.Name)
	}
}
