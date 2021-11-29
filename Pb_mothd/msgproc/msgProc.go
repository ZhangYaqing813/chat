package msgproc

import (
	msg "chat/Message_type"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Messager struct {
	Conn net.Conn
}

// Msgjson 消息的序列化，返回一个字符串
func (M *Messager) Msgjson(message msg.Messages) (b []byte) {
	b, err := json.Marshal(message)
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
func (M *Messager) MsgReader(b []byte) (messages msg.Messages) {
	var buf []byte
	_, err := M.Conn.Read(b[:4])
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
func (M *Messager) Transmit(messages msg.Messages) {

}
