package wx

import (
	"fmt"
)

type Plugin interface {
	Name() string
	Type() int
	String() string
}

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

//  管理器负责 管理插件 和 命中插件执行。
type PluginManager interface {
	//  获取插件
	Get(name string) Plugin
	//  添加插件
	Add(plugin Plugin)
	//  支持 String 输出
	String() string

	GetAll(t int) map[string]Plugin

	Run(wrapper Plugin, sess *Session, msg *ReceivedMessage)
}

//  type=1
type PluginText interface {
	HandleText(sess *Session, msg *ReceivedMessage)
}

//  type=51
type PluginInit interface {
	HandleInit(sess *Session, msg *ReceivedMessage)
}

//  type = 49
type PluginLink interface {
	HandleLink(sess *Session, msg *ReceivedMessage)
}

type PluginWrapper struct {
	plugin  Plugin
	enabled bool
}

type defaultPluginManager struct {
	// all
	plugins map[string]Plugin

	//  mapper
	users  map[string]Plugin
	groups map[string]Plugin
	all    map[string]Plugin
}

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

func (dmp *defaultPluginManager) String() string {
	return fmt.Sprintf("dmp count:%d all:%d users:%d groups:%d\n", len(dmp.plugins), len(dmp.users), len(dmp.groups), len(dmp.all))
}

func (dmp *defaultPluginManager) GetAll(t int) map[string]Plugin {
	if t&PluginTypeAll == PluginTypeAll {
		return dmp.all
	}

	if t&PluginTypeUser == PluginTypeUser {
		return dmp.users
	}

	if t&PluginTypeGroup == PluginTypeGroup {
		return dmp.groups
	}

	return nil
}

func (dmp *defaultPluginManager) Run(plugin Plugin, sess *Session, msg *ReceivedMessage) {

	switch msg.MsgType {
	case MSG_TEXT:
		//  1 文本
		if p, ok := plugin.(PluginText); ok {
			p.HandleText(sess, msg)
			return
		}
	case MSG_INIT:
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
