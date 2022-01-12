package cfgm

var defaultContext = NewConfigManageContext(&ConfigManageContextOptions{
	CommandLinePrefix:    "",
	ConfigFilePathPrefix: "",
})

func Register(path string, ptrToConfigObject interface{}, callback ConfigManageCallback) {
	defaultContext.Register(path, ptrToConfigObject, callback)
}

func Get(path string, ptr interface{}) bool {
	return defaultContext.Get(path, ptr)
}

func Init() []error {
	return defaultContext.Init()
}
