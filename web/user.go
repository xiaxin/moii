package web

// User WEB 用户接口
type User interface {
	GetUID() uint
	GetUserName() string
	GetToken() string
}

type user struct {
	Id       uint
	Username string
	Token    string
}

func (u user) GetUID() uint {
	return u.Id
}

func (u user) GetUserName() string {
	return u.Username
}

func (u user) GetToken() string {
	return ""
}
func NewUser(id uint, username, token string) User {
	return &user{
		Id:       id,
		Username: username,
		Token:    token,
	}
}

func NewNilUser() User {
	return &user{
		Id:       0,
		Username: "访客",
	}
}
