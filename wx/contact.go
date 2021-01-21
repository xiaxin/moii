package wx

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/xiaxin/moii/log"
)

// Contact 联系人
type Contact struct {
	user    *User
	cm      ContactManager
	plugins []Plugin
}

func NewContact(user *User) *Contact {
	contact := &Contact{
		user: user,
	}

	if contact.IsGroup() && nil == contact.cm {
		contact.cm = NewContactManager(nil)
		contact.cm.AddUsers(contact.user.MemberList)
	}

	return contact
}

func (c *Contact) IsGroup() bool {
	return c.user.MemberCount > 0
}

func (c *Contact) GetUser() *User {
	return c.user
}

func (c *Contact) GetUin() int32 {
	return c.user.Uin
}

func (c *Contact) GetUsername() string {
	return c.user.UserName
}

func (c *Contact) GetNickname() string {
	nickname := c.user.NickName
	return nickname
}

func (c *Contact) GetRemarkName() string {
	return c.user.RemarkName
}

func (c *Contact) GetSignature() string {
	return c.user.Signature
}

func (c *Contact) GetDisplayName() string {
	// TODO  parent
	return c.user.DisplayName
}

func (c *Contact) GetMember(username string) *Contact {

	if c.IsGroup() {
		return c.cm.Get(username)
	}
	return nil
}

func (c *Contact) GetContactManager() ContactManager {
	return c.cm
}

func (c *Contact) GetMemberByNickname(nickname string) *Contact {
	if c.IsGroup() {
		return c.cm.GetByNickname(nickname)
	}
	return nil
}

type contactManager struct {
	mu      sync.Mutex
	data    map[string]*Contact
	session *Session
	//  联系人管理 默认是否安装插件
	plugin bool
}

func NewContactManager(sess *Session) ContactManager {
	cm := &contactManager{
		session: sess,
		data:    make(map[string]*Contact),
		plugin:  true,
	}

	if nil == sess {
		cm.plugin = false
	}

	return cm
}

func (cm *contactManager) AddUsers(users []*User) {
	var groups []*User

	for _, user := range users {
		if !strings.Contains(user.UserName, "@@") {
			cm.AddUser(user.UserName, user)
		} else {
			groups = append(groups, user)
		}
	}

	for _, user := range cm.GetGroupContact(groups) {
		cm.AddGroup(user.UserName, user)
	}
}

//  增加普通用户
func (cm *contactManager) AddUser(username string, user *User) {
	cm.addUser(username, user, false)
}

func (cm *contactManager) AddGroup(username string, user *User) {
	cm.addUser(username, user, true)
}

func (cm *contactManager) addUser(username string, user *User, group bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	contact := NewContact(user)

	if true == cm.plugin {
		//  添加插件
		for _, v := range cm.session.PluginManager.GetAll(PluginTypeAll) {
			contact.plugins = append(contact.plugins, v)
		}

		if contact.IsGroup() {
			for _, v := range cm.session.PluginManager.GetAll(PluginTypeGroup) {
				contact.plugins = append(contact.plugins, v)
			}
		} else {
			for _, v := range cm.session.PluginManager.GetAll(PluginTypeUser) {
				contact.plugins = append(contact.plugins, v)
			}
		}
	}

	cm.data[username] = contact

}

func (cm *contactManager) Get(username string) *Contact {
	if user, ok := cm.data[username]; ok {
		return user
	}
	return nil
}

func (cm *contactManager) GetUser(username string) *User {
	if contact := cm.Get(username); nil != contact {
		return contact.user
	}
	return nil
}

func (cm *contactManager) GetByNickname(nickname string) *Contact {

	for _, v := range cm.data {

		//log.Infof("[@] group.user %s", v.user)
		//log.Infof("[@] group.nickname:%s nickname:%s", v.user.NickName, nickname)
		if v.user.NickName == nickname || v.user.DisplayName == nickname {
			return v
		}
	}
	return nil
}

func (cm *contactManager) GetNickName(username string) (string, error) {
	user, ok := cm.data[username]

	if !ok {
		return "", errors.New(fmt.Sprintf("Cm GetNickName not found %s", username))
	}
	return user.GetNickname(), nil
}

func (cm *contactManager) GetContact(username string) *Contact {
	return cm.Get(username)
}

func (cm *contactManager) GetData() map[string]*Contact {
	return cm.data
}

func (cm *contactManager) PullBatchContractByUsername(username string) {
	var users []*User

	cm.PullBatchContracts(append(users, &User{
		EncryChatRoomID: "",
		UserName:        username,
	}))
}

func (cm *contactManager) PullBatchContracts(user []*User) {
	for _, user := range cm.GetGroupContact(user) {
		if !strings.Contains(user.UserName, "@@") {
			cm.AddUser(user.UserName, user)
		} else {
			cm.AddGroup(user.UserName, user)
		}
	}
}

// TODO 依赖宿主SESSION
func (cm *contactManager) GetGroupContact(groups []*User) []*User {
	//  联系人中的群组管理 session = nil
	if nil == cm.session {
		return nil
	}

	b, err := WebWxBatchGetContact(cm.session.WxConfig, cm.session.WxXMLConfig, cm.session.Cookies, groups)

	batch, err := ParseInitResponse(b)
	if nil != err {
		log.Error("ParseInitResponse Error %s", err)
		return nil
	}

	return batch.ContactList
}
