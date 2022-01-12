package cfgm

import (
	"errors"
	"fmt"
	"github.com/SnowPhoenix0105/cfgm/internal/check"
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
)

type ConfigManageContextOptions struct {
	CommandLinePrefix    string
	ConfigFilePathPrefix string
}

type ConfigManageCallback func(err error) error

type registerItem struct {
	Path     string
	Obj      interface{}
	Callback ConfigManageCallback
}

type ConfigManageContext struct {
	root          *tree.Node
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
		root:          nil,
		options:       options,
		registerItems: nil,
	}
}

func (ctx *ConfigManageContext) Register(path string, ptrToConfigObject interface{}, callback ConfigManageCallback) {
	if !check.IsPtr(ptrToConfigObject) {
		panic(errors.New(fmt.Sprintf("cfgm register error: Config Object is not registered with it's pointer")))
	}
	ctx.registerItems = append(ctx.registerItems, registerItem{
		Path:     path,
		Obj:      ptrToConfigObject,
		Callback: callback,
	})
}

func (ctx *ConfigManageContext) Get(path string, ptr interface{}) bool {
	panic("not implement")
}

func (ctx *ConfigManageContext) Init() []error {
	panic("not implement")
}
