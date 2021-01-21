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
	MsgText = 1 // text message

	MsgImg    = 3  // image message
	MSG_VOICE = 34 // voice message
	MsgVoice  = 35

	MsgFV       = 37 // friend verification message
	MsgPF       = 40 // POSSIBLEFRIEND_MSG
	MsgSCC      = 42 // shared contact card
	MSG_VIDEO   = 43 // video message
	MSG_EMOTION = 47 // gif
	MsgEmotion  = 47

	MSG_LOCATION = 48 // location message
	MsgLocation  = 48

	MsgLink = 49 // shared link message

	MSG_VOIP = 50 // VOIPMSG
	MsgVoip  = 50

	MsgInit = 51 // wechat init message

	MSG_VOIPNOTIFY = 52 // VOIPNOTIFY
	MsgVoipnotify  = 52

	MSG_VOIPINVITE = 53 // VOIPINVITE
	MsgVoipinvite  = 53

	MSG_SHORT_VIDEO = 62 // short video message
	MsgShortVideo   = 62

	MSG_SYSNOTICE = 9999 // SYSNOTICE
	MsgSysNotice  = 9999

	MSG_SYS = 10000 // system message
	MsgSys  = 10000

	MSG_WITHDRAW = 10002 // withdraw notification message
	MsgWithDraw  = 10002
)

// User 用户
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

	ChatRoomID int
	KeyWord    string

	MemberStatus int

	EncryChatRoomID string
	IsOwner         int
}

func (u *User) String() string {
	return fmt.Sprintf("[user] username:%s nickname:%s, remark:%s, mc:%d sign:%s", u.UserName, u.NickName, u.RemarkName, u.MemberCount, u.Signature)
}

// SyncKeyList 用于同步消息的KEY
type SyncKeyList struct {
	Count int
	List  []SyncKey
}

// SyncKey KEY|VAL
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

// ReceivedMessage 收到的消息
type ReceivedMessage struct {
	IsGroup       bool
	MsgID         string
	Content       string
	FromUserName  string
	ToUserName    string
	Who           string
	MsgType       int
	SubType       int
	OriginContent string
	At            string
	AtUser        *User
	URL           string

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

// BaseRequest Response 中 BaseRequest
type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

// InitRequest TODO
type InitRequest struct {
	BaseRequest        *BaseRequest
	Msg                interface{}
	SyncKey            *SyncKeyList
	RR                 int
	Code               int
	FromUserName       string
	ToUserName         string
	ClientMsgID        int
	ClientMediaID      int
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

// BaseResponse 响应结构体
type BaseResponse struct {
	Ret    int
	ErrMsg string
}

// InitResponse TODO1
type InitResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	Count        int           `json:"Count"`
	ContactList  []*User       `json:"ContactList"`
	Skey         string        `json:"Skey"`
	SyncKey      *SyncKeyList  `json:"SyncKey"`
	User         *User
}

// ContactResponse TODO
type ContactResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
	MemberCount  int           `json:"MemberCount"`
	MemberList   []*User       `json:"MemberList"`
	Seq          int           `json:"Seq"`
}

// ReceiveResponse TODO
type ReceiveResponse struct {
	BaseResponse    *BaseResponse     `json:"BaseResponse"`
	AddMsgCount     int               `json:"AddMsgCount"`
	AddMsgList      []*ReceiveMessage `json:"AddMsgList"`
	ModContactCount int               `json:"ModContactCount"`
	ModContactList  []*User           `json:"ModContactList"`
}

// ReceiveMessage TODO
type ReceiveMessage struct {
	MsgID                string
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
	MediaID              string
	URL                  string
	AddMsgType           int
	StatusNotifyCode     int
	StatusNotifyUserName string
	RecommendInfo        *RecommendInfo
	ForwardFlag          int
	AppInfo              *AppInfo
	HasProductID         int
	Ticket               string
	ImgHeight            int
	ImgWidth             int
	SubMsgType           int
	NewMsgID             uint64
	OriContent           string
	EncryFileName        string
}

// OplogRequest TODO
type OplogRequest struct {
	BaseRequest *BaseRequest
	UserName    string
	CmdID       int
	RemarkName  string
}

// OplogResponse TODO
type OplogResponse struct {
	BaseResponse *BaseResponse `json:"BaseResponse"`
}
