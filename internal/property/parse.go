package property

import (
	"errors"
	"fmt"
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

func enterNode(ptr *node, path string) *node {
	next, ok := ptr.sub[path]
	if !ok {
		next = newNode()
		ptr.sub[path] = next
	}
	return next
}

func appendProperty(root *node, prop string) error {
	beg := 0
	ptr := root
	length := len(prop)
	for end := 1; end < length; {
		if prop[end] == '.' {
			ptr = enterNode(ptr, prop[beg:end])
			beg = end + 1
			end += 2
			continue
		}
		if prop[end] == '=' {
			ptr = enterNode(ptr, prop[beg:end])
			if len(ptr.value) != 0 {
				// this no has been set by another property
				return errors.New(fmt.Sprintf("property conflict at %s", prop[0:end]))
			}
			if end == length-1 {
				// set this node as an empty node
				return nil
			}
			ptr.value = prop[end+1 : length]
			return nil
		}
		end++
	}
	if beg == length {
		return errors.New(fmt.Sprintf("invalid property: %s", prop))
	}
	ptr = enterNode(ptr, prop[beg:length])
	// set this node as an empty node
	return nil
}

func ParseFromPropertyList(propList []string) (Record, error) {
	root := newNode()
	for _, prop := range propList {
		err := appendProperty(root, prop)
		if err != nil {
			return Record{root: nil}, err
		}
	}
	return Record{root: root}, nil
}

func FixTree(record Record, root *tree.Node, time tree.ModifyTime) error {
	env := &fixEnv{walker: tree.WriteFrom(root, time)}
	env.fixNode(record.root)
	return nil
}
