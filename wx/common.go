package wx

import (
	"fmt"
	"strconv"
	"strings"
)

/**
- common
- message
- request
*/

const (
	// msg types
	MSG_TEXT        = 1     // text message
	MSG_IMG         = 3     // image message
	MSG_VOICE       = 34    // voice message
	MSG_FV          = 37    // friend verification message
	MSG_PF          = 40    // POSSIBLEFRIEND_MSG
	MSG_SCC         = 42    // shared contact card
	MSG_VIDEO       = 43    // video message
	MSG_EMOTION     = 47    // gif
	MSG_LOCATION    = 48    // location message
	MsgLink         = 49    // shared link message
	MSG_VOIP        = 50    // VOIPMSG
	MSG_INIT        = 51    // wechat init message
	MSG_VOIPNOTIFY  = 52    // VOIPNOTIFY
	MSG_VOIPINVITE  = 53    // VOIPINVITE
	MSG_SHORT_VIDEO = 62    // short video message
	MSG_SYSNOTICE   = 9999  // SYSNOTICE
	MSG_SYS         = 10000 // system message
	MSG_WITHDRAW    = 10002 // withdraw notification message
)

type User struct {
	Uin               int32
	UserName          string
	NickName          string
	HeadImgUpdateFlag int
	ContactType       int
	ContactFlag       int

	// GROUP
	MemberCount int
	MemberList  []*User

	RemarkName       string
	HideInputBarFlag int
	Sex              int
	Signature        string
	VerifyFlag       int

	OwnerUin int

	PYQuanPin       string
	PYInitial       string
	RemarkPYInitial string
	RemarkPYQuanPin string

	StarFriend     int
	AppAccountFlag int
	Statues        int

	AttrStatus int

	Province string
	City     string

	Alias string

	SnsFlag int

	UniFriend   int
	DisplayName string

	ChatRoomId int
	KeyWord    string

	MemberStatus int

	EncryChatRoomId string
	IsOwner         int
}

func (u *User) String() string {
	return fmt.Sprintf("[user] username:%s nickname:%s, remark:%s, mc:%d sign:%s", u.UserName, u.NickName, u.RemarkName, u.MemberCount, u.Signature)
}

type SyncKeyList struct {
	Count int
	List  []SyncKey
}

type SyncKey struct {
	Key int
	Val int
}

func (s *SyncKeyList) String() string {
	strs := make([]string, 0)
	for _, v := range s.List {
		strs = append(strs, strconv.Itoa(v.Key)+"_"+strconv.Itoa(v.Val))
	}
	return strings.Join(strs, "|")
}

type VerifyUser struct {
	Value            string
	VerifyUserTicket string
}

type RecommendInfo struct {
	UserName   string
	NickName   string
	QQNum      int
	Province   string
	City       string
	Content    string
	Signature  string
	Alias      string
	Scene      int
	VerifyFlag int
	AttrStatus uint32
	Sex        int
	Ticket     string
	OpCode     int
}

type AppInfo struct {
	AppID string
	Type  int
}

type ReceivedMessage struct {
	IsGroup       bool
	MsgId         string
	Content       string
	FromUserName  string
	ToUserName    string
	Who           string
	MsgType       int
	SubType       int
	OriginContent string
	At            string
	AtUser        *User
	Url           string

	RecommendInfo *RecommendInfo
}

func (msg *ReceivedMessage) String() string {
	return ""
}

// TextMessage: text message struct
type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
}

// MediaMessage
type MediaMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
	MediaId      string
}

// EmotionMessage: gif/emoji message struct
type EmotionMessage struct {
	ClientMsgId  int
	EmojiFlag    int
	FromUserName string
	LocalID      int
	MediaId      string
	ToUserName   string
	Type         int
}

// Request
type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

type InitRequest struct {
	BaseRequest        *BaseRequest
	Msg                interface{}
	SyncKey            *SyncKeyList
	RR                 int
	Code               int
	FromUserName       string
	ToUserName         string
	ClientMsgId        int
	ClientMediaId      int
	TotalLen           int
	StartPos           int
	DataLen            int
	MediaType          int
	Scene              int
	Count              int
	List               []*User
	Opcode             int
	SceneList          []int
	SceneListCount     int
	VerifyContent      string
	VerifyUserList     []*VerifyUser
	VerifyUserListSize int
	skey               string
	MemberCount        int
	MemberList         []*User
	Topic              string
}

// Response 响应结构体
type BaseResponse struct {
	Ret    int
	ErrMsg string
}

type InitResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	Count        int           `json:"Count"`
	ContactList  []*User       `json:"ContactList"`
	Skey         string        `json:"Skey"`
	SyncKey      *SyncKeyList  `json:"SyncKey"`
	User         *User
}

type ContactResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	MemberCount  int           `json:"MemberCount"`
	MemberList   []*User       `json:"MemberList"`
	Seq          int           `json:"Seq"`
}

type ReceiveResponse struct {
	BaseResponse    *BaseResponse     `json:"BaseResponse"`
	AddMsgCount     int               `json:"AddMsgCount"`
	AddMsgList      []*ReceiveMessage `json:"AddMsgList"`
	ModContactCount int               `json:"ModContactCount"`
	ModContactList  []*User           `json:"ModContactList"`
}

type ReceiveMessage struct {
	MsgId                string
	FromUserName         string
	ToUserName           string
	MsgType              int
	Content              string
	Status               int
	ImgStatus            int
	CreateTime           int
	VoiceLength          int
	PlayLength           int
	FileName             string
	FileSize             string
	MediaId              string
	Url                  string
	AddMsgType           int
	StatusNotifyCode     int
	StatusNotifyUserName string
	RecommendInfo        *RecommendInfo
	ForwardFlag          int
	AppInfo              *AppInfo
	HasProductId         int
	Ticket               string
	ImgHeight            int
	ImgWidth             int
	SubMsgType           int
	NewMsgId             uint64
	OriContent           string
	EncryFileName        string
}

type OplogRequest struct {
	BaseRequest *BaseRequest
	UserName    string
	CmdId       int
	RemarkName  string
}

type OplogResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
}
