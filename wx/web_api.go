package wx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultClient = NewHttpClient()
	LogoutSign    = make(map[int]int)
)

type Header map[string]string

func init() {
	LogoutSign[1100] = 1
	LogoutSign[1101] = 1
	LogoutSign[1102] = 1
	LogoutSign[1205] = 1
}

func WebJsLogin(common *WebConfig) (string, error) {

	km := url.Values{}
	km.Add("appid", common.AppId)
	km.Add("fun", "new")
	km.Add("lang", common.Lang)
	km.Add("redirect_uri", common.RedirectUri)
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri := common.LoginUrl + "/jslogin?" + km.Encode()

	body, err := DefaultClient.Get(uri, nil)

	if nil != err {
		return "", fmt.Errorf("WebApi.JsLogin error: %s", err)
	}

	ss := strings.Split(string(body), "\"")
	if len(ss) < 2 {
		return "", fmt.Errorf("jslogin response invalid, %s", string(body))
	}
	return ss[1], nil
}

func WebNewLoginPage(common *WebConfig, xc *WebXmlConfig, uri string) ([]*http.Cookie, error) {
	u, _ := url.Parse(uri)
	km := u.Query()
	km.Add("fun", "new")
	uri = common.CgiUrl + "/webwxnewloginpage?" + km.Encode()
	resp, err := DefaultClient.FetchReponse("GET", uri, []byte(""), Header{})

	if nil != err {
		return nil, fmt.Errorf("WebApi.WebNewLoginPage error: %s", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err := xml.Unmarshal(body, xc); err != nil {
		return nil, err
	}
	if xc.Ret != 0 {
		return nil, fmt.Errorf("xc.Ret != 0, %s", string(body))
	}
	return resp.Cookies(), nil
}

// TODO 针对返回值 进行合理性优化
func WebLogin(common *WebConfig, uuid, tip string) (string, error) {
	km := url.Values{}
	km.Add("tip", tip)
	km.Add("uuid", uuid)
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("loginicon", "true")
	uri := common.LoginUrl + "/cgi-bin/mmwebwx-bin/login?" + km.Encode()
	body, err := DefaultClient.Get(uri, nil)

	if nil != err {
		return "", fmt.Errorf("WebApi.Login error: %s", err)
	}

	strb := string(body)
	if strings.Contains(strb, "window.code=200") &&
		strings.Contains(strb, "window.redirect_uri") {
		ss := strings.Split(strb, "\"")
		if len(ss) < 2 {
			return "", fmt.Errorf("parse redirect_uri fail, %s", strb)
		}
		return ss[1], nil
	}

	return "", fmt.Errorf("login response, %s", strb)
}

func WebWxInit(comm *WebConfig, ce *WebXmlConfig) ([]byte, error) {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("skey", ce.Skey)
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

	uri := comm.CgiUrl + "/webwxinit?" + km.Encode()

	js := InitRequest{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			comm.DeviceID,
		},
	}

	b, _ := json.Marshal(js)

	body, err := DefaultClient.PostJsonByte(uri, b)

	if nil != err {
		return nil, fmt.Errorf("WebApi.WebWxInit Post Request error: %s", err)
	}

	return body, nil
}

func WebWxStatusNotify(config *WebConfig, ce *WebXmlConfig, bot *User) (int, error) {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("lang", config.Lang)

	uri := config.CgiUrl + "/webwxstatusnotify?" + km.Encode()

	js := InitRequest{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			config.DeviceID,
		},
		Code:         3,
		FromUserName: bot.UserName,
		ToUserName:   bot.UserName,
		ClientMsgId:  int(time.Now().Unix()),
	}

	b, _ := json.Marshal(js)

	body, err := DefaultClient.PostJsonByte(uri, b)

	if nil != err {
		return -1, fmt.Errorf("WebApi.WebWxStatusNotify PostRequest Error:%s", err)
	}

	response, err := ParseInitResponse(body)

	if nil != err {
		return -1, fmt.Errorf("WebApi.WebWxStatusNotify ParseResponse Error:%s", err)
	}

	return response.BaseResponse.Ret, nil
}

func WebWxGetContact(config *WebConfig, ce *WebXmlConfig, cookies []*http.Cookie) ([]byte, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("seq", "0")
	km.Add("skey", ce.Skey)
	uri := config.CgiUrl + "/webwxgetcontact?" + km.Encode()

	js := InitRequest{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			config.DeviceID,
		},
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	DefaultClient.SetJar(jar)
	body, _ := DefaultClient.PostJsonByte(uri, b)
	return body, nil
}

func SyncCheck(comm *WebConfig, ce *WebXmlConfig, cookies []*http.Cookie, server string, skl *SyncKeyList) (retcode int, selector int, err error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
	km.Add("sid", ce.Wxsid)
	km.Add("uin", ce.Wxuin)
	km.Add("skey", ce.Skey)
	km.Add("deviceid", comm.DeviceID)
	km.Add("synckey", skl.String())
	km.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))
	uri := "https://" + server + "/cgi-bin/mmwebwx-bin/synccheck?" + km.Encode()

	js := InitRequest{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			comm.DeviceID,
		},
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	DefaultClient.SetJar(jar)

	body, err := DefaultClient.GetByte(uri, b)

	if nil != err {
		return -1, -1, fmt.Errorf("WebApi.SyncCheck error: %s", err)
	}

	strb := string(body)
	reg := regexp.MustCompile("window.synccheck={retcode:\"(\\d+)\",selector:\"(\\d+)\"}")
	sub := reg.FindStringSubmatch(strb)
	retcode = 0
	selector = 0
	if len(sub) >= 2 {
		retcode, _ = strconv.Atoi(sub[1])
		selector, _ = strconv.Atoi(sub[2])
	}

	return retcode, selector, nil
}

func WebWxSync(comm *WebConfig, ce *WebXmlConfig, cookies []*http.Cookie, msg chan []byte, skl *SyncKeyList) ([]*http.Cookie, error) {

	km := url.Values{}
	km.Add("skey", ce.Skey)
	km.Add("sid", ce.Wxsid)
	km.Add("lang", comm.Lang)
	km.Add("pass_ticket", ce.PassTicket)

	uri := comm.CgiUrl + "/webwxsync?" + km.Encode()

	js := InitRequest{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			comm.DeviceID,
		},
		SyncKey: skl,
		RR:      ^int(time.Now().Unix()) + 1,
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	DefaultClient.SetJar(jar)

	resp, body, _ := DefaultClient.PostJsonByteForResp(uri, b)

	response, err := ParseInitResponse(body)

	if err != nil {
		return nil, err
	}

	retcode := response.BaseResponse.Ret

	if retcode != 0 {
		//  TODO -3003
		return nil, fmt.Errorf("BaseResponse.Ret %d BaseResponse.ErrMsg %d", retcode, response.BaseResponse.Ret)
	}

	msg <- body

	skl.List = skl.List[:0]
	skl1 := response.SyncKey
	skl.Count = skl1.Count
	skl.List = append(skl.List, skl1.List...)

	return resp.Cookies(), nil
}

func WebWxBatchGetContact(comm *WebConfig, ce *WebXmlConfig, cookies []*http.Cookie, cl []*User) ([]byte, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("type", "ex")
	uri := comm.CgiUrl + "/webwxbatchgetcontact?" + km.Encode()

	js := InitRequest{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			comm.DeviceID,
		},
		Count: len(cl),
		List:  cl,
	}

	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("WebApi.WebWxBatchGetContact error: %s", err)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", comm.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func GetLoginAvatar(resp string) (string, error) {
	match := regexp.MustCompile(`window.userAvatar = '(.+)'`).
		FindStringSubmatch(resp)
	if len(match) < 2 {
		return "", errors.New("login avatar not found")
	}
	return match[1], nil
}

func WebWxSendMsg(comm *WebConfig, ce *WebXmlConfig, cookies []*http.Cookie,
	from, to string, msg string) ([]byte, error) {

	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)

	uri := comm.CgiUrl + "/webwxsendmsg?" + km.Encode()

	js := InitRequest{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			comm.DeviceID,
		},
		Msg: &TextMessage{
			Type:         1,
			Content:      msg,
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  int(time.Now().Unix() * 1e4),
		},
	}

	b, _ := json.Marshal(js)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	DefaultClient.SetJar(jar)
	body, _ := DefaultClient.PostJsonByte(uri, b)
	return body, nil
}

func WebWxOplog(conf *WebConfig, ce *WebXmlConfig, cookies []*http.Cookie, user *User, name string) ([]byte, error) {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)

	uri := conf.CgiUrl + "/webwxoplog?" + km.Encode()

	js := &OplogRequest{
		UserName:   user.UserName,
		CmdId:      2,
		RemarkName: name,
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			conf.DeviceID,
		},
	}

	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("WebApi.WebWxOplog error: %s", err)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", conf.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
