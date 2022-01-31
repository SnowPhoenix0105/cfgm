package controller

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
)

type ConfigManageCallback func(err error) error

type registerItem struct {
	Path     []string
	Obj      interface{}
	Callback ConfigManageCallback
	Error    error
}

type ConfigManageContext struct {
	root          *tree.Node
	configObject  map[string]interface{}
	options       *ConfigManageContextOptions
	registerItems []registerItem
}

func NewConfigManageContext(options *ConfigManageContextOptions) *ConfigManageContext {
	if len(options.CommandLinePrefix) == 0 {
		options.CommandLinePrefix = "-D"
	}
	if len(options.ConfigFilePathPrefix) == 0 {
		options.ConfigFilePathPrefix = "--config="
	}
	return &ConfigManageContext{
		root:          tree.NewNode(),
		options:       options,
		registerItems: nil,
	}
}

func (ctx *ConfigManageContext) Get(path string, ptr interface{}) bool {
	panic("not implement")
}
