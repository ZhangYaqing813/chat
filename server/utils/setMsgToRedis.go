package utils

import (
	messagetype "chat/Message_type"
	chatlog "chat/chatLog"
	redism "chat/server/dba/redis"
)

// 最终调用在信息转发的方法下

/*
	消息缓存需求：
		1、消息转发完成后需要对一下发送者，接收者 进行缓存，一条消息要缓存多份。
			发送者缓存到发送的列表，接受者缓存到接收的列表。
		2、不在线的需要缓存到对应用户的未读消息记录
		3、未读消息需要用户上线后有server 端进行转发，然后删除对应用户的未读消息记录，移动到接收的消息库

		redis hash  key 格式：
			userAReFromUserB 	该user 接收到的信息。
	        userASendToUserB  	该user 发送的信息
	        userAUnreadMessage 	该user 未读的信息
*/

// SetMsgToRedis sendUser:消息发送者；receiveUser：消息接收者列表
//messages: 完整的信息；field：二级key，使用时间戳
func SetMsgToRedis(sendUser string, receiveUser []string, messages messagetype.Messages, field string) {
	// 确认接收用户是否在线
	for _, user := range receiveUser {
		for _, v := range OnlineUsers {
			if user != v {
				key := user + "UnreadMessage"
				redism.MyRedis.AddMessage(key, messages, field)
			} else {
				// 缓存发送者发送过的消息
				key := sendUser + "SendTo" + user
				redism.MyRedis.AddMessage(key, messages, field)
				// 缓存消息到接受者
				key = user + "ReFrom" + sendUser
				redism.MyRedis.AddMessage(key, messages, field)
			}
		}
	}
}

// GetMsgFromRedis 二次封装 redis  GetMessage
/*	需要进行数据的处理,区分获取信息的位置
	获取redis 中的信息分为两种：
	1、或去某个用户的未读信息
		只需要用户名即可
	2、获取历史信息
		2.1、获取自己发送信息
			可以控制根据发送对象进而获取不同的自己发送的信息
		2.2、获取自己收到的信息
			可以控制根据发送对象进而获取不同的自己收到的信息
	本次不考虑这些功能，只需要取出历史信息，转发给对应的客户端即可

*/

// GetUnReadMsgRedis 此时的设计，field 为空，取key 中所有的值
func GetUnReadMsgRedis(userName string) (str []string) {
	key := userName + "UnreadMessage"
	str, err := redism.MyRedis.GetMessage(key)
	if err != nil {
		chatlog.Std.Errorf("GetUnReadMsgRedis get unreadmsg failed err = %v", err)
	}
	return
}

func GetMsgFromRedis(key string, field []string) (str []string) {

	return
}
