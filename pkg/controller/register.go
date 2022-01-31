package controller

import (
	"errors"
	"fmt"
	"github.com/SnowPhoenix0105/cfgm/internal/check"
	"strings"
)

func (ctx *ConfigManageContext) Register(path string, ptrToConfigObject interface{}, callback ConfigManageCallback) {
	if !check.IsPtr(ptrToConfigObject) {
		panic(errors.New(fmt.Sprintf("cfgm register error: Config Object is not registered with it's pointer")))
	}
	ctx.registerItems = append(ctx.registerItems, registerItem{
		Path:     strings.Split(path, "."),
		Obj:      ptrToConfigObject,
		Callback: callback,
		Error:    nil,
	})
}
