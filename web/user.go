package web

type User interface {
	GetUid() uint
	GetUserName() string
	GetToken() string
}

