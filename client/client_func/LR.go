package client_func

import (
	messagetype "chat/Message_type"
	"chat/Pb_mothd/msgproc"
	"encoding/json"
	"fmt"
)

type LR struct {
	//conn net.Conn
	msgproc.Messager
}

func (L *LR) Login(loginmsg messagetype.LoginMsg) (code int, error string) {
	var msg messagetype.Messages
	var recmsg messagetype.LResMsg
	//组装要发送的信息到msg, 消息类型是 messagetype.LoginMsgType
	msg.Type = messagetype.LoginMsgType
	msg.Data = string(L.Msgjson(loginmsg))

	//调用 msgsender 发送组装好的数据
	L.MsgSender(msg)
	//接受server 返回的数据
	msg = L.MsgReader()
	// 解析返回的数据
	err := json.Unmarshal([]byte(msg.Data), &recmsg)
	if err != nil {
		fmt.Println("解析login response message failed", err)
		return
	}
	if recmsg.Code == 200 {
		// 用一个协成 使client 和server 保持通讯
		go L.keepSession()
	}
	//返回解析完的数据
	return recmsg.Code, recmsg.Error
}

func (L *LR) Register(register_message messagetype.RegMsg) (code int, error string) {
	var msg messagetype.Messages
	var recmsg messagetype.LResMsg

	// 组装注册信息
	msg.Type = messagetype.ResMsg
	msg.Data = string(L.Msgjson(register_message))

	// 发送组装完成后的信息
	L.MsgSender(msg)

	//接收server 返回的数据信息
	err := json.Unmarshal([]byte(L.MsgReader().Data), &recmsg)
	if err != nil {
		fmt.Println("register response message unmarshal failed ", err)
		return
	}
	//返回解析完成后的数据
	return recmsg.Code, recmsg.Error
}

// Chat 发送对话信息
func (L *LR) Chat(dialogue messagetype.Dialogue) {

	var message messagetype.Messages
	//设置要发送的消息头尾 messagetype.ChatMode
	message.Type = messagetype.ChatMode
	// 格式化要发送的信息内用
	message.Data = string(L.Msgjson(dialogue))
	//发送消息
	L.MsgSender(message)

}

// 保持和server 端的连接
func (L *LR) keepSession() {
	newmsg := L.MsgReader()
	fmt.Println(newmsg)
}
