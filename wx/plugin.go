package wx

import (
	"fmt"
)

const (
	//  user plugin
	PluginTypeUser = 1 << iota
	//  group plugin
	PluginTypeGroup
	PluginTypeAll = PluginTypeUser | PluginTypeGroup
)

var (
	DefaultPluginManager PluginManager
)

func init() {
	DefaultPluginManager = NewDefaultPluginManager()
}

// type PluginWrapper struct {
// 	plugin  Plugin
// 	enabled bool
// }

type defaultPluginManager struct {
	// all
	plugins map[string]Plugin

	//  mapper
	users  map[string]Plugin
	groups map[string]Plugin
	all    map[string]Plugin
}

// NewDefaultPluginManager 创建一个插件管理器
func NewDefaultPluginManager() PluginManager {
	return &defaultPluginManager{
		plugins: make(map[string]Plugin),
		users:   make(map[string]Plugin),
		groups:  make(map[string]Plugin),
		all:     make(map[string]Plugin),
	}
}

func (dpm *defaultPluginManager) Get(name string) Plugin {
	if plugin, ok := dpm.plugins[name]; ok {
		return plugin
	}

	return nil
}

func (dpm *defaultPluginManager) Add(plugin Plugin) {
	name := plugin.Name()

	dpm.plugins[name] = plugin

	if plugin.Type()&PluginTypeAll == PluginTypeAll {
		dpm.all[name] = dpm.plugins[name]
	} else if plugin.Type()&PluginTypeUser == PluginTypeUser {
		dpm.users[name] = dpm.plugins[name]
	} else if plugin.Type()&PluginTypeGroup == PluginTypeGroup {
		dpm.groups[name] = dpm.plugins[name]
	}
}

// String 文字说明
func (dpm *defaultPluginManager) String() string {
	return fmt.Sprintf("dmp count:%d all:%d users:%d groups:%d\n", len(dpm.plugins), len(dpm.users), len(dpm.groups), len(dpm.all))
}

// GetAll 获取所有
func (dpm *defaultPluginManager) GetAll(t int) map[string]Plugin {
	if t&PluginTypeAll == PluginTypeAll {
		return dpm.all
	}

	if t&PluginTypeUser == PluginTypeUser {
		return dpm.users
	}

	if t&PluginTypeGroup == PluginTypeGroup {
		return dpm.groups
	}

	return nil
}

// Run 运行
func (dpm *defaultPluginManager) Run(plugin Plugin, sess *Session, msg *ReceivedMessage) {

	switch msg.MsgType {
	case MsgText:
		//  1 文本
		if p, ok := plugin.(PluginText); ok {
			p.HandleText(sess, msg)
			return
		}
	case MsgInit:
		//  51
		if p, ok := plugin.(PluginInit); ok {
			p.HandleInit(sess, msg)
			return
		}
	case MsgLink:
		//
	}

	//  todo 兜底

}
