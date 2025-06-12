package plugins

import "github.com/yek-j/filemanager/config"

type Plugin interface {
	Process(cfg *config.Config) error
	GetName() string
	GetDescription() string
}