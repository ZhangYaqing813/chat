package message_type

import (
	"net"
)

const (
	LoginMsgType = "LoginMsg"
	RegMsgType   = "RegMsgType"
	ResMsg       = "ResMsg"
	ExitMsgType  = "Exit"
	ChatMode     = "Chat"
	UPDATE       = "UPDATE"
)

// status code
const (
	SUCCESS = 200
	FAILED  = 500
	// SINGLE 单聊模式
	SINGLE = 1
	//MULTIPLE 多人聊天模式
	MULTIPLE = 2
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

// UserOnline  获取redis 中在线用户的用户名
type UserOnline struct {
	UsersOnline []string `json:"usersonline,[]string"`
}

// UserOlineIntoRedis 用于写入redis
type UserOlineIntoRedis struct {
	UserName string   `json:"username"`
	UserConn net.Conn `json:"userconn,net.Conn"`
}

// Dialogue 对话信息

type Dialogue struct {
	// sendmod 可以为 S/s 单聊，M/m 多聊
	SendMod int `json:"sendmod"`
	// 信息发送的对象，
	ToUsers []string `json:"users,[]string"`
	//消息的发送者
	Sender string `json:"sender"`
	//消息内容
	Content string `json:"content"`
	//发送消息的时间
	SendTime string `json:"sendtime"`
}
