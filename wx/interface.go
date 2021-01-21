package wx

// ContactManager 联系人管理器
type ContactManager interface {
	AddUsers(users []*User)
	AddUser(username string, user *User)
	AddGroup(username string, user *User)
	Get(username string) *Contact
	GetByNickname(nickname string) *Contact
	GetContact(username string) *Contact
	GetUser(username string) *User
	GetData() map[string]*Contact
	PullBatchContractByUsername(username string)
	PullBatchContracts(user []*User)
	GetGroupContact(groups []*User) []*User
}

// PluginManager 管理器负责 管理插件 和 命中插件执行。
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

// Plugin 插件
type Plugin interface {
	Name() string
	Type() int
	String() string
}

// PluginText type = 1
type PluginText interface {
	HandleText(sess *Session, msg *ReceivedMessage)
}

// PluginInit type = 51
type PluginInit interface {
	HandleInit(sess *Session, msg *ReceivedMessage)
}

// PluginLink type = 49
type PluginLink interface {
	HandleLink(sess *Session, msg *ReceivedMessage)
}
