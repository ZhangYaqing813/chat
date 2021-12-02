package message_type

const (
	LoginMsgType = "LoginMsg"
	RegMsgType   = "RegMsgType"
	ResMsg       = "ResMsg"
	ExitMsgType  = "Exit"
)

//LoginMsg 登录是向server 提交的用户信息
type LoginMsg struct { // 用户注册消息的结构体
	// 用户Id
	UserID int `json:"userid" DB:"userid"`
	//用户名
	UserName string `json:"username" DB:"username"`
	//用户密码
	Password string `json:"password" DB:"password"`
}

// LResMsg 用户登录状态信息
type LResMsg struct {
	//登录状态码
	Code int `json:"code"`
	// 登录错误时的错误信息
	Error string `json:"Error"`
}

//Messages 要发送的信息
type Messages struct {
	//发送的消息类型
	Type string `json:"type"`
	// 发送信息的内容
	Data string `json:"data"`
}

type RegMsg struct {
	// 用户Id
	UserID int `json:"userid" DB:"userid"`
	//用户密码
	Password string `json:"password" DB:"password"`
	//用户名
	UserName string `json:"username" DB:"username"`
	//UserEmail
	UserEmail string `json:"email" DB:"email"`
	//Sex
	Sex string `json:"sex" DB:"sex"`
}

type UserOnlie struct {
	//
	Type       string `json:"type"`
	UsersOnlie string `json:"usersonlie"`
}

// status code

const (
	SUCCESS = 200
	FAILED  = 500
)
