package msgproc

import (
	msg "chat/Message_type"
	chatlog "chat/chatLog"
	"chat/server/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
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
func (M *Messager) MsgReader() (messages msg.Messages, err error) {
	var buf [8192]byte
	_, err = M.Conn.Read(buf[:4])

	if err != nil {
		//fmt.Println("接收客户端数据失败", err)
		chatlog.Std.Errorf("Read msg lenth failed %v", err)
		return messages, err
	}
	n := binary.BigEndian.Uint32(buf[:4])
	_, err = M.Conn.Read(buf[:n])
	if err != nil {
		//fmt.Println("接收客户端数据失败", err)
		chatlog.Std.Errorf("Read msg failed %v", err)
		return messages, err
	}
	data := M.UnJson(buf[:n])
	return data, err

}

// MsgSender 发送信息
// 1、获取message 消息的长度，装入 buf 的前四位
// 2、发送buf[:4] 给server 端，告知server端本次发送数据的长度
// 3、
func (M *Messager) MsgSender(messages msg.Messages) (err error) {
	var buf [8192]byte
	// 计算出要发送的消息体的长度，然后发送个对端，根据发送的长度进行数据校验
	//可以在消息体message 中增加一个字段，接收者收到后取出进行对比
	binary.BigEndian.PutUint32(buf[:4], uint32(len(M.Msgjson(messages))))
	n, err := M.Conn.Write(buf[:4])
	if err != nil || n != 4 {
		chatlog.Std.Errorf("Send msg lenth failed %v", err)
		return err
	} else {
		_, err = M.Conn.Write(M.Msgjson(messages))
		if err != nil {
			chatlog.Std.Errorf("Send msg failed %v", err)
			return err
		}
	}
	return err
}

// Transmit 转发信息
// 消息转发确定的几个问题
/*
1、如何调用conn.write() 把信息转发到对应的user
	1.1 conn的信息直接以参数的形式进行传参，修改 msgSender 方法
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
		fmt.Println("sendToUser", sendToUser)

		// 增加一个判断逻辑，接收消息的用户必须在线，如果不在线需要将消息缓存到该用户的未读信息中
		if len(sendToUser) > 0 && utils.GetUserStatus(sendToUser, utils.OnlineUsers) {
			M.Conn = utils.OnlineUserInfo[sendToUser]
			err := M.MsgSender(messages)
			if err != nil {
				chatlog.Std.Error(err)
			}
		} else {
			// 把该用户添加到NotOnline 列表中
			utils.SetNotOnLine(sendToUser)
			// 添加到Redis UnReadMessage
			utils.SetMsgToRedis(sendToUser, dialogueMessage.ToUsers, messages, dialogueMessage.SendTime)
			continue
		}
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
	if online == true {
		message.Data = "user " + userName + " online"
	} else {
		message.Data = "user " + userName + " offline"
	}
	// 发送更新消息给在线用户列表
	for _, user := range utils.OnlineUsers {
		if user != "" {
			M.Conn = utils.OnlineUserInfo[user]
			chatlog.Std.Infof("%s 连接地址：%v", user, utils.OnlineUserInfo[user])
			err := M.MsgSender(message)
			if err != nil {
				chatlog.Std.Errorf("NotifyOnline MsgSender err= %v", err)
			}
		} else {
			continue
		}
	}
}
