package web

// User WEB 用户接口
type User interface {
	GetUID() uint
	GetUserName() string
	GetToken() string
}
