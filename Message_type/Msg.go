package message_type

import "net"

const (
	LoginMsgType = "LoginMsg"
	RegMsgType   = "RegMsgType"
	ResMsg       = "ResMsg"
	ExitMsgType  = "Exit"
	ChatMode     = "Chat"
	UPDATE       = "UPDATE"
	GETSENDMSG   = "SelfSend"
	GETRECMSG    = "SelfReceive"
	UNREADMSG    = "UnReadMsg"
	RESPONSETF   = "ResponseTf"
	RESPONSE     = "Response"
)

// status code
const (
	SUCCESS = 200
	FAILED  = 500
	// SINGLE 单聊模式
	SINGLE = 1
	//MULTIPLE 多人聊天模式
	MULTIPLE = 2
	// USERONLINE 用户在线
	USERONLINE = 1
	// USEROFFLINE 用户不在线
	USEROFFLINE = 2
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

// Response 增加一个新的结构类型，用于回复client 请求后的消息体
type Response struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

//Messages 要发送的信息
type Messages struct {
	//发送的消息类型
	Type string `json:"type"`
	// 发送信息的内容
	Data string `json:"data"`
}

// RegMsg 注册信息
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

//
type UserUpdate struct {
	// 更新字段的字段名 如：email 、password 等
	FieldName string `json:"fieldName"`
	// 更新的内用
	NewContent string `json:"newContent"`
}

// UserOnline  获取redis 中在线用户的用户名
type UserOnline struct {
	UsersOnline []string `json:"usersonline,[]string"`
}

// UserOlineIntoRedis 用于写入redis
// 暂时未用
type UserOlineIntoRedis struct {
	UserName string   `json:"username"`
	UserConn net.Conn `json:"userconn"`
}

// UserOlineOutRedis 用于写入redis
// 暂时未用
type UserOlineOutRedis struct {
	UserName string   `json:"username"`
	UserConn net.Conn `json:"userconn"`
}

// Dialogue 对话信息

type Dialogue struct {
	// sendMod 可以为 S/s 单聊，M/m 多聊
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

type History struct {
	//对谁的聊天记录
	MsgType  string   `json:"msgtype"`
	User     string   `json:"user"`
	Messages []string `json:"messages"`
}

// UnMsg 未读消息体
type UnMsg struct {
	Data string `json:"data"`
}

// ResponseTf 转发回复
type ResponseTf struct {
	Data string `json:"data"`
}

// ClientConfig 客户端配置文件
type ClientConfig struct {
	ConnectIP string `ini:"connectIp"`
	Port      int    `ini:"port"`
}

// RedisConfig redis 配置文件
type RedisConfig struct {
	Host string `ini:"host"`
	Port string `ini:"rPort"`
	Auth string `ini:"auth"`
	Db   int    `ini:"db"`
}

// MysqlConfig 配置文件
type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     string `ini:"mPort"`
	Username string `ini:"username"`
	Pwd      string `ini:"pwd"`
	Dbname   string `ini:"dbname"`
}

// SvcConfig Svc 服务器监听配置
type SvcConfig struct {
	ListenIP string `ini:"listenIp"`
	Port     int    `ini:"serverPort"`
}

// ServerConfig  服务器配置文件
type ServerConfig struct {
	RedisConfig `ini:"redis""`
	MysqlConfig `ini:"mysql"`
	SvcConfig   `ini:"server"`
}
