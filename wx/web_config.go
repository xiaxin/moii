package wx

import "encoding/xml"

var (
	// DefaultWebConfig TODO: 字段变更为何时字段
	DefaultWebConfig = &WxConfig{
		AppID:      "wx782c26e4c19acffb",
		LoginURL:   "https://login.weixin.qq.com",
		Lang:       "zh_CN",
		DeviceID:   "e" + GetRandomStringFromNum(15),
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		SyncSrv:    "webpush.wx.qq.com",
		UploadURL:  "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
		MediaCount: 0,
	}
	// URLPool TODO: 此处数据 如何更新
	URLPool = []WebURLGroup{
		{"wx2.qq.com", "file.wx2.qq.com", "webpush.wx2.qq.com"},
		{"wx8.qq.com", "file.wx8.qq.com", "webpush.wx8.qq.com"},
		{"qq.com", "file.wx.qq.com", "webpush.wx.qq.com"},
		{"web2.wechat.com", "file.web2.wechat.com", "webpush.web2.wechat.com"},
		{"wechat.com", "file.web.wechat.com", "webpush.web.wechat.com"},
	}
)

type WxConfig struct {
	AppID       string
	LoginURL    string
	Lang        string
	DeviceID    string
	UserAgent   string
	CgiURL      string
	CgiDomain   string
	SyncSrv     string
	UploadURL   string
	MediaCount  uint32
	RedirectURI string
}

// WebXMLConfig 登录配置
type WebXMLConfig struct {
	XMLName     xml.Name `xml:"error"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale int      `xml:"isgrayscale"`
}

// WebURLGroup TODO
type WebURLGroup struct {
	IndexURL  string
	UploadURL string
	SyncURL   string
}
