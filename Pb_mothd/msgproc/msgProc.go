package msgproc

import (
	msg "chat/Message_type"
	"chat/server/untils"
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
	//fmt.Println("Msgjson ",v)
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("message marshal failed")
	}
	return b
}

// UnJson 反序列化，返回一个msg
func (M *Messager) UnJson(b []byte) (payload msg.Messages) {
	err := json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println("反序列化失败", err)
	}
	return payload
}

//MsgReader 读取信息
func (M *Messager) MsgReader() (messages msg.Messages) {
	var buf [8192]byte
	_, err := M.Conn.Read(buf[:4])
	if err != nil {
		fmt.Println("d读取消息体长度失败", err)
		return
	}
	n := binary.BigEndian.Uint32(buf[:4])
	_, err = M.Conn.Read(buf[:n])
	if err != nil {
		fmt.Println("读取消息失败", err)
		return
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
		fmt.Println("消息长度发送失败", err)
		return
	} else {
		_, err = M.Conn.Write(M.Msgjson(messages))
		if err != nil {
			fmt.Println("消息发送失败", err)
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

	msg_send_user := untils.GetOlineUser(dialogueMessage.ToUsers)
	for _, user := range msg_send_user {
		M.Conn = user.UserConn
		M.MsgSender(messages)
	}

	// 2、 封装 需要转发的message ,

	// 3、

}
