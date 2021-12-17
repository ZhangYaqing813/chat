package msgproc

import (
	msg "chat/Message_type"
	chatlog "chat/chatLog"
	"chat/server/untils"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
)

type Messager struct {
	Conn net.Conn
}

// Msgjson 消息的序列化，返回一个字符串
//修改返回值类型
func (M *Messager) Msgjson(v interface{}) (b []byte) {
	b, err := json.Marshal(v)
	if err != nil {
		chatlog.Std.Fatalf("Marshal message failed  %v", err)
	}
	return b
}

// UnJson 反序列化，返回一个msg
func (M *Messager) UnJson(b []byte) (payload msg.Messages) {
	err := json.Unmarshal(b, &payload)
	if err != nil {
		chatlog.Std.Fatalf("UnMarshal message failed  %v", err)
	}
	return payload
}

//MsgReader 读取信息
func (M *Messager) MsgReader() (messages msg.Messages) {
	var buf [8192]byte
	_, err := M.Conn.Read(buf[:4])

	if err == io.EOF || err != nil {
		chatlog.Std.Errorf("Read msg lenth failed %v", err)
	}
	n := binary.BigEndian.Uint32(buf[:4])
	_, err = M.Conn.Read(buf[:n])
	if err != nil {
		chatlog.Std.Errorf("Read msg failed %v", err)
	}
	data := M.UnJson(buf[:n])
	return data

}

// MsgSender 发送信息
// 1、获取message 消息的长度，装入 buf 的前四位
// 2、发送buf[:4] 给server 端，告知server端本次发送数据的长度
// 3、
func (M *Messager) MsgSender(messages msg.Messages) {
	var buf [8192]byte
	// 计算出要发送的消息体的长度，然后发送个对端，根据发送的长度进行数据校验
	//可以在消息体message 中增加一个字段，接收者收到后取出进行对比
	binary.BigEndian.PutUint32(buf[:4], uint32(len(M.Msgjson(messages))))
	n, err := M.Conn.Write(buf[:4])
	if err != nil || n != 4 {
		chatlog.Std.Fatalf("Send msg lenth failed %v", err)
		return
	} else {
		_, err = M.Conn.Write(M.Msgjson(messages))
		if err != nil {
			chatlog.Std.Fatalf("Send msg failed %v", err)
			return
		}
	}
}

// Transmit 转发信息
// 消息转发确定的几个问题
/*
1、如何调用conn.write() 把信息转发到对应的user
	1.1 conn的信息直接以参数的形式进行传参，修改 msgsender 方法
	1.2 修改G.conn 值，然后进行方法调用，问题：会不会影响其他的正常消息发送
2、消息中携带的信息应该有哪些
	发送者
	接受者
	消息内容
	时间
3、消息转发完成后如何写入redis
*/
func (M *Messager) Transmit(dialogueMessage msg.Dialogue, messages msg.Messages) {
	// 1、 根据 dialogueMessage.ChatSignal.SendMod  模式获取发送消息对象的内存地址
	for _, sendToUser := range dialogueMessage.ToUsers {
		M.Conn = untils.OnlineUserInfo[sendToUser]
		M.MsgSender(messages)
	}
}

//NotifyOnline 通知用户
// NotifyOnline online 值为 true || false
//users []string 当前在线用户李彪
//username 刚上线的用户名

func (M *Messager) NotifyOnline(userName string, online bool) {
	var message msg.Messages
	message.Type = msg.UPDATE

	// 确定通知的是上线还是下线
	if online {
		message.Data = userName + "online"
	} else {
		message.Data = userName + "offline"
	}
	// 发送更新消息给在线用户列表

	for _, user := range untils.OnlineUsers {

		M.Conn = untils.OnlineUserInfo[user]
		chatlog.Std.Infof("%s 连接地址：%v", user, untils.OnlineUserInfo[user])
		M.MsgSender(message)
	}
}
