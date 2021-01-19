package wx

import "encoding/xml"

var (
	// TODO 字段变更为何时字段
	DefaultWebConfig = &WebConfig{
		AppId:      "wx782c26e4c19acffb",
		LoginUrl:   "https://login.weixin.qq.com",
		Lang:       "zh_CN",
		DeviceID:   "e" + GetRandomStringFromNum(15),
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		SyncSrv:    "webpush.wx.qq.com",
		UploadUrl:  "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
		MediaCount: 0,
	}
	//  TODO
	URLPool = []WebUrlGroup{
		{"wx2.qq.com", "file.wx2.qq.com", "webpush.wx2.qq.com"},
		{"wx8.qq.com", "file.wx8.qq.com", "webpush.wx8.qq.com"},
		{"qq.com", "file.wx.qq.com", "webpush.wx.qq.com"},
		{"web2.wechat.com", "file.web2.wechat.com", "webpush.web2.wechat.com"},
		{"wechat.com", "file.web.wechat.com", "webpush.web.wechat.com"},
	}
)

type WebConfig struct {
	AppId       string
	LoginUrl    string
	Lang        string
	DeviceID    string
	UserAgent   string
	CgiUrl      string
	CgiDomain   string
	SyncSrv     string
	UploadUrl   string
	MediaCount  uint32
	RedirectUri string
}

type WebXmlConfig struct {
	XMLName     xml.Name `xml:"error"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale int      `xml:"isgrayscale"`
}

type WebUrlGroup struct {
	IndexUrl  string
	UploadUrl string
	SyncUrl   string
}
