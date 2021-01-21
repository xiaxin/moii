package wx

import (
	"errors"
	"fmt"

	"github.com/mdp/qrterminal"

	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/xiaxin/moii/log"
)

const (
	ErrorUserLogout = "user logout"

	// TerminalMode 命令行
	TerminalMode = 0
	// TerminalModeGoland GoLand
	TerminalModeGoland = 1
)

// Session 会话周期
type Session struct {
	WxConfig       *WxConfig
	WxXMLConfig    *WebXMLConfig
	muCookie       sync.RWMutex
	Cookies        []*http.Cookie
	SynKeyList     *SyncKeyList
	ContactManager ContactManager
	QrcodePath     string //qrcode path
	QrcodeUUID     string //uuid
	PluginManager  PluginManager
	CreateTime     int64
	LastMsgID      string
	//Api             *wxweb.ApiV2
	OnLoginAvatar func(string) error
	AfterLogin    func() error
	Owner         *User
	qrmode        int
	State         uint8

	// 退出 Chan
	Exit chan struct{}
}

// CreateSession 创建SESSION qrmode
func CreateSession(qrmode int) (*Session, error) {
	session := &Session{
		WxConfig:    DefaultWebConfig, // 配置
		WxXMLConfig: &WebXMLConfig{},  // XML配置 TODO
		//QrcodeUUID:  uuid,
		// TODO 迁移完成时 需要 使用 Api接入方式 处理
		//Api:         api,
		CreateTime: time.Now().Unix(), // 时间

		PluginManager: DefaultPluginManager, //  插件管理器

		OnLoginAvatar: func(string) error { // TODO
			return nil
		},
		AfterLogin: func() error { // TODO
			return nil
		},
		qrmode: qrmode,
		State:  0,
		Exit:   make(chan struct{}),
	}

	//  创建联系人管理器
	session.ContactManager = NewContactManager(session)

	return session, nil
}

// Login wx.JsLogin 获取UUID
func (s *Session) Login() error {
	uuid, err := WebJsLogin(s.WxConfig)

	log.Infof("[session] [login] uuid:%s", uuid)

	if err != nil {
		return err
	}

	s.QrcodeUUID = uuid

	if s.qrmode == TerminalMode {
		qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)
	} else if s.qrmode == TerminalModeGoland {
		config := qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: qrterminal.WHITE,
			WhiteChar: qrterminal.BLACK,
			QuietZone: 1,
		}
		qrterminal.GenerateWithConfig("https://login.weixin.qq.com/l/"+uuid, config)
	}

	return nil
}

//func (s *Session) SetPipeline(pipeline pipeline.Pipeline) {
//	s.Pipeline = pipeline
//}

// SetCookies 设置Cookie
func (s *Session) SetCookies(cookies []*http.Cookie) {
	s.muCookie.Lock()
	defer s.muCookie.Unlock()
	s.Cookies = cookies
}

// GetCookies 获取Cookie
func (s *Session) GetCookies() []*http.Cookie {
	s.muCookie.RLock()
	defer s.muCookie.RUnlock()
	return s.Cookies
}

// Producer 消息生产者
func (s *Session) Producer(msg chan []byte, errChan chan error) {
	log.Info("entering synccheck loop")
loop1:
	for {
		var (
			ret int
			sel int // selector
			err error
		)
		for i := 0; i <= 10; i++ {
			//  TODO 检查状态
			ret, sel, err = SyncCheck(s.WxConfig, s.WxXMLConfig, s.GetCookies(), s.WxConfig.SyncSrv, s.SynKeyList)
			if err != nil {
				if i >= 10 {
					log.Errorf("SyncCheck error:(%s) try time %d ret %d selector %d\n", err.Error(), i, ret, sel)
				} else {
					log.Infof("SyncCheck uin(%d) time %d ret %d selector %d\n", s.Owner.Uin, i, ret, sel)
				}
			} else {
				break
			}
		}

		// 正常
		if ret == 0 {
			// check success
			// new message
			for i := 0; i <= 10; i++ {
				cookies, err := WebWxSync(s.WxConfig, s.WxXMLConfig, s.GetCookies(), msg, s.SynKeyList)
				if err != nil {
					if i >= 10 {
						log.Error("Err WebWxSync try  %s try %d", err.Error(), i)
					} else {
						log.Info("WebWxSync uin %d tiem %d", s.Owner.Uin, i)
					}
				} else {
					if cookies != nil {
						s.SetCookies(cookies)
					}
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		} else if s.IsLogout(ret) { //1100 失败/登出微信
			errChan <- fmt.Errorf(ErrorUserLogout)
			break loop1
		} else {
			errChan <- fmt.Errorf("unhandled exception ret %d", ret)
			break loop1
		}
	}

}

// IsLogout 是否退出
func (s *Session) IsLogout(code int) bool {
	_, has := LogoutSign[code]
	return has
}

// ScanWaiter （检查二维码状态）
func (s *Session) ScanWaiter(onAvatar func(string) error) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	//  每秒发起一次检查
	for range ticker.C {
		redirectURI, err := WebLogin(s.WxConfig, s.QrcodeUUID, "0")
		if err != nil {
			//  TODO
			log.Warnf("ScanWaiter Error %s", err)
			switch {
			case strings.Contains(err.Error(), "window.code=400"),
				strings.Contains(err.Error(), "window.code=500"),
				strings.Contains(err.Error(), "window.code=0"):
				return err
			case strings.Contains(err.Error(), "window.code=201"):
				// 二维码被扫描
				avatar, err := GetLoginAvatar(err.Error())
				if err != nil {
					return err
				}
				if err := onAvatar(avatar); err != nil {
					return err
				}
			}
		} else {
			s.WxConfig.RedirectURI = redirectURI
			s.AnalizeVersion(s.WxConfig.RedirectURI)
			break
		}
	}
	return nil
}

// AnalizeVersion TODO
func (s *Session) AnalizeVersion(uri string) {
	u, _ := url.Parse(uri)

	// version may change
	s.WxConfig.CgiDomain = u.Scheme + "://" + u.Host
	s.WxConfig.CgiURL = s.WxConfig.CgiDomain + "/cgi-bin/mmwebwx-bin"

	for _, urlGroup := range URLPool {
		if strings.Contains(u.Host, urlGroup.IndexURL) {
			s.WxConfig.SyncSrv = urlGroup.SyncURL
			s.WxConfig.UploadURL = fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json", urlGroup.UploadURL)
			return
		}
	}
}

// LoginAndServe 处理逻辑
func (s *Session) LoginAndServe(useCache bool) error {

	var (
		err error
	)

	//  @TODO  如果 微信App 主动退出网页登录。这里还会继续登录。
	if !useCache {
		if s.GetCookies() != nil {
			// confirmWaiter
		}

		if err := s.ScanWaiter(s.OnLoginAvatar); err != nil {
			return err
		}
		var cookies []*http.Cookie
		// update cookies
		if cookies, err = WebNewLoginPage(s.WxConfig, s.WxXMLConfig, s.WxConfig.RedirectURI); err != nil {
			return err
		}
		s.SetCookies(cookies)

	}

	//  初始化 ContactList 中 有常用 用户 和 群组
	jb, err := WebWxInit(s.WxConfig, s.WxXMLConfig)

	if err != nil {
		return err
	}

	//  解析 INIT 数据
	msg, err := ParseInitResponse(jb)

	if nil != err {
		return err
	}

	if s.IsLogout(msg.BaseResponse.Ret) {
		return fmt.Errorf(ErrorUserLogout)
	}

	//  添加常用联系人
	s.ContactManager.AddUsers(msg.ContactList)

	//  设置 机器人 用户
	s.Owner = msg.User
	s.SynKeyList = msg.SyncKey

	//  状态通知
	ret, err := WebWxStatusNotify(s.WxConfig, s.WxXMLConfig, s.Owner)
	if err != nil {
		return err
	}
	if ret != 0 {
		return fmt.Errorf("WebWxStatusNotify fail, %d", ret)
	}

	//  获取个人联系人
	cb, err := WebWxGetContact(s.WxConfig, s.WxXMLConfig, s.GetCookies())
	if err != nil {
		return err
	}

	//  解析联系人信息
	contact, err := ParseContactResponse(cb)
	if nil != err {
		return err
	}
	s.ContactManager.AddUsers(contact.MemberList)

	s.AfterLogin()

	if err := s.serve(); err != nil {
		//s.Pipeline.Stop()
		return err
	}
	return nil
}

func (s *Session) serve() error {

	s.State = 1 // 已启动

	msg := make(chan []byte, 1000)
	// syncheck
	errChan := make(chan error)
	//  开始处理消息
	go s.Producer(msg, errChan)

	for {
		select {
		case m := <-msg:
			resp, err := ParseReceiveResponse(m)

			if nil != err {
				log.Infof("[session] [serve] [parse-receive] error:%s", err)
			} else {
				go s.Consumer(resp)
			}
		case err := <-errChan:
			return err
		}
	}
}

// Consumer  TODO 分析具体的消息体 整理出消息协议。
func (s *Session) Consumer(receive *ReceiveResponse) {

	//  人员变更时都会触发 ModContactList 和 ModContactCount 的变更
	// receive, err := ParseReceiveResponse(msg)

	// if nil != err {
	// 	log.Error("ParseReceiveResponse Error %s", err)
	// 	return
	// }

	//  TODO 更新动作 应该新起一个  GoRoutine 执行，并且在MsgList 之后执行
	if receive.ModContactCount > 0 {
		for _, v := range receive.ModContactList {
			if v.MemberCount > 0 {
				log.Infof("MODC [群组更新] ID: %s Nickname: %s MemberCount:%d\n", v.UserName, v.NickName, v.MemberCount)
				s.ContactManager.AddGroup(v.UserName, v)
			} else {
				log.Infof("MODC [联系人更新] ID:%s Nickname: %s RemarkName: %s\n", v.UserName, v.NickName, v.RemarkName)
				s.ContactManager.AddUser(v.UserName, v)
			}
		}
	}

	//  新消息
	if receive.AddMsgCount > 0 {
		for _, v := range receive.AddMsgList {
			if nil == s.ContactManager.GetContact(v.FromUserName) {
				s.ContactManager.PullBatchContractByUsername(v.FromUserName)
			}

			//  过滤机器人自己说的话
			if v.FromUserName == s.Owner.UserName {
				continue
			}

			// 分析
			rmsg := s.Analize(v)

			contact := s.ContactManager.GetContact(v.FromUserName)

			if nil != contact {
				go s.Recv(contact, rmsg)
			} else {
				log.Error("contact:%s is not found", v.FromUserName)
			}
		}
	}
}

// Analize TODO
func (s *Session) Analize(msg *ReceiveMessage) *ReceivedMessage {
	rmsg := &ReceivedMessage{
		MsgID:         msg.MsgID,
		OriginContent: msg.Content,
		FromUserName:  msg.FromUserName,
		ToUserName:    msg.ToUserName,
		MsgType:       msg.MsgType,
		SubType:       msg.SubMsgType,
		URL:           msg.URL,
	}

	// friend verify message
	if rmsg.MsgType == MsgFV {
		rmsg.RecommendInfo = msg.RecommendInfo
	}

	if strings.Contains(rmsg.FromUserName, "@@") {
		rmsg.IsGroup = true

		//  群组信息
		ss := strings.Split(rmsg.OriginContent, ":<br/>")
		if len(ss) > 1 {
			rmsg.Who = ss[0]
			rmsg.Content = ss[1]
		} else {
			rmsg.Who = s.Owner.UserName
			rmsg.Content = rmsg.OriginContent
		}

	} else {
		// none group message
		rmsg.Who = rmsg.FromUserName
		rmsg.Content = rmsg.OriginContent
	}

	if rmsg.MsgType == MsgText &&
		//  检查内容
		len(rmsg.Content) > 1 &&
		//  检查 @
		strings.HasPrefix(rmsg.Content, "@") {
		// @someone
		ss := strings.Split(rmsg.Content, "\u2005")

		if len(ss) == 2 {
			rmsg.At = ss[0] + "\u2005"
			rmsg.Content = ss[1]

			//  TODO 通过Nickname查询User
			from := s.ContactManager.Get(msg.FromUserName)
			if nil == from {
				//log.Warnf("[@] username not found %s", msg.FromUserName)
			} else {

				//log.Infof("[@] nickname:%s len:%d", ss[0][1:], len(rmsg.At))
				member := from.GetMemberByNickname(ss[0][1:])

				if nil != member {
					rmsg.AtUser = member.GetUser()
				}
			}
		}
	}
	return rmsg
}

// Recv 接收消息
func (s *Session) Recv(c *Contact, msg *ReceivedMessage) {
	// TODO 这里要处理 panic
	for _, plugin := range c.plugins {
		go s.PluginManager.Run(plugin, s, msg)
	}
}

// UpdateRemarkName 更新别名
func (s *Session) UpdateRemarkName(username string, name string) error {

	user := s.ContactManager.GetUser(username)

	if nil == user {
		return errors.New("update remarkname error: user is nil")
	}

	ret, err := WebWxOplog(s.WxConfig, s.WxXMLConfig, s.Cookies, user, name)

	if nil != err {
		return err
	}

	log.Infof("update result:%s\n", string(ret))

	return nil
}

// AddPlugin 添加插件
func (s *Session) AddPlugin(plugin Plugin) {
	s.PluginManager.Add(plugin)
}
