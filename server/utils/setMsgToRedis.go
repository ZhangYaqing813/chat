package utils

import (
	messagetype "chat/Message_type"
	redism "chat/server/dba/redis"
)

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
//需要进行数据的处理
//区分获取信息的位置

func GetMsgFromRedis(field []string) (str []string) {

	return
}
