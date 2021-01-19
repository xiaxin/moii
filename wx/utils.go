package wx

import (
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func unmarshal(msg []byte, v interface{}) error {
	err := json.Unmarshal(msg, v)
	return err
}

//  解析消息体
func ParseReceiveResponse(msg []byte) (*ReceiveResponse, error) {

	message := new(ReceiveResponse)

	if err := unmarshal(msg, message); nil != err {
		return nil, err
	}

	return message, nil
}

//  解析初始化信息
func ParseInitResponse(msg []byte) (*InitResponse, error) {
	init := new(InitResponse)
	err := unmarshal(msg, init)
	return init, err
}

func ParseContactResponse(msg []byte) (*ContactResponse, error) {
	contact := new(ContactResponse)
	err := unmarshal(msg, contact)
	return contact, err
}

func GetRandomStringFromNum(length int) string {
	bytes := []byte("0123456789")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RealTargetUserName(session *Session, msg *ReceivedMessage) string {
	if session.Owner.UserName == msg.FromUserName {
		return msg.ToUserName
	}
	return msg.FromUserName
}
