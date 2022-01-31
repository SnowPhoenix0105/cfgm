package property

import (
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"strings"
)

type CmdPropertyParseOptions struct {
	ConfigFilePrefix string
	PropertyPrefix   string
}

func ParseFromCmd(cmdLines []string, options *CmdPropertyParseOptions) (string, Record, error) {
	propList := make([]string, 0)
	configFilePath := ""
	for _, cmd := range cmdLines {
		if strings.HasPrefix(cmd, options.ConfigFilePrefix) {
			configFilePath = strings.TrimPrefix(cmd, options.ConfigFilePrefix)
		} else if strings.HasPrefix(cmd, options.PropertyPrefix) {
			propList = append(propList, strings.TrimPrefix(cmd, options.PropertyPrefix))
		}
	}
	record, err := ParseFromPropertyList(propList)
	return configFilePath, record, err
}

func ParseFromPropertyList(propList []string) (Record, error) {
	// TODO
	panic("not implement")
}

func FixTree(record Record, root *tree.Node, time tree.ModifyTime) error {
	// TODO
	panic("not implement")
}
