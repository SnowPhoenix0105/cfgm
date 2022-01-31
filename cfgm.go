package cfgm

import (
	"github.com/SnowPhoenix0105/cfgm/pkg/controller"
)

var defaultContext = controller.NewConfigManageContext(&controller.ConfigManageContextOptions{
	CommandLinePrefix:    "",
	ConfigFilePathPrefix: "",
})

func Register(path string, ptrToConfigObject interface{}, callback controller.ConfigManageCallback) {
	defaultContext.Register(path, ptrToConfigObject, callback)
}

func Get(path string, ptr interface{}) bool {
	return defaultContext.Get(path, ptr)
}

func Init() []error {
	return defaultContext.Init()
}
