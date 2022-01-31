package controller

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/SnowPhoenix0105/cfgm/internal/json2tree"
	"github.com/SnowPhoenix0105/cfgm/internal/obj2tree"
	"github.com/SnowPhoenix0105/cfgm/internal/property"
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"github.com/SnowPhoenix0105/cfgm/internal/tree2obj"
	"os"
	"strings"
)

const (
	modifyTimeInvalid tree.ModifyTime = iota
	modifyTimeBuild
	modifyTimeMerge
	modifyTimeEnv
	modifyTimeCmd
	modifyTimeAssign
)

func (ctx *ConfigManageContext) addConfigObjectToTree(walker tree.Walker, configObject interface{}, path []string) error {
	for _, p := range path {
		walker.EnterObj(p)
	}
	err := obj2tree.AppendTo(configObject, walker)
	for range path {
		walker.Exit()
	}
	return err
}

func (ctx *ConfigManageContext) buildTreeFromObjectConfig() bool {
	ok := true
	walker := tree.WriteFrom(ctx.root, modifyTimeBuild)
	for i, item := range ctx.registerItems {
		err := ctx.addConfigObjectToTree(walker, item.Obj, item.Path)
		if err != nil {
			ok = false
			ctx.registerItems[i].Error = err
		}
	}
	return ok
}

func (ctx *ConfigManageContext) mergeTreeByFileConfig(filePath string) error {
	if strings.HasSuffix(filePath, ".json") {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		reader := bufio.NewReader(file)
		return json2tree.Merge(ctx.root, reader, modifyTimeMerge)
	}
	return errors.New(fmt.Sprintf("unsupported file type :%s", filePath))
}

func (ctx *ConfigManageContext) parseCmd() (string, property.Record, error) {
	return property.ParseFromCmd(os.Args[1:], &property.CmdPropertyParseOptions{
		ConfigFilePrefix: "--config=",
		PropertyPrefix:   "-D",
	})
}

func (ctx *ConfigManageContext) fixTree(record property.Record) error {
	return property.FixTree(record, ctx.root, modifyTimeCmd)
}

func resolveItem(root *tree.Node, item *registerItem, ch chan<- error) {
	if item.Error != nil {
		ch <- item.Callback(item.Error)
		return
	}
	walker := tree.ReadFrom(root)
	for _, p := range item.Path {
		if !walker.TryEnterObj(p) {
			if DEBUG {
				ch <- item.Callback(nil)
				panic("TryEnterObj() fail with key from path")
			}
		}
	}
	tree2obj.RefillFrom(walker, item.Obj, modifyTimeBuild)
	ch <- item.Callback(nil)
}

func (ctx *ConfigManageContext) invokeCallbacks(e error) []error {
	ch := make(chan error, 1)
	if e != nil {
		for i := range ctx.registerItems {
			ctx.registerItems[i].Error = e
		}
	}

	// because resolveItem() only read the AST, so it's ok to invoke it in parallel
	for i := range ctx.registerItems {
		go resolveItem(ctx.root, &ctx.registerItems[i], ch)
	}

	// join all goroutines and build the result
	errs := make(map[error]struct{}, 0)
	for range ctx.registerItems {
		err := <-ch
		if err == nil {
			continue
		}
		errs[err] = struct{}{}
	}
	ret := make([]error, 0)
	for err := range errs {
		ret = append(ret, err)
	}
	return ret
}

func (ctx *ConfigManageContext) Init() []error {
	ok := ctx.buildTreeFromObjectConfig()
	if !ok {
		return ctx.invokeCallbacks(nil)
	}
	filePath, record, err := ctx.parseCmd()
	if err != nil {
		return ctx.invokeCallbacks(err)
	}
	if len(filePath) != 0 {
		err = ctx.mergeTreeByFileConfig(filePath)
		if err != nil {
			return ctx.invokeCallbacks(err)
		}
	}
	err = ctx.fixTree(record)
	if err != nil {
		return ctx.invokeCallbacks(err)
	}
	return ctx.invokeCallbacks(nil)
}
