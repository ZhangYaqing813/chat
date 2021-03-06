package utils

import (
	chatlog "chat/chatLog"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

/*
维护三个表的意义在于：
	OnlineUserInfo 维护了整个用户的信息，用户名和通讯地址
	OnlineUsers 只用于维护在线用户，主要用户在线用户列表更新，减少一次 OnlineUserInfo 遍历
	notOnline   当消息转发的时候，需要检测下用户是否在线，如果没有在线怎把要转发消息存入该用户的未读消息库
*/

var OnlineUserInfo = make(map[string]net.Conn)
var OnlineUsers = make([]string, 8)
var NotOnline = make([]string, 16)

// AddUser 向 OnlineUser 中添加一个上线的用户
func AddUser(userName string, conn net.Conn) {
	//向OnlineUserInfo 中添加user 的conn 信息
	OnlineUserInfo[userName] = conn

	//向OnlineUsers 中添加user
	OnlineUsers = append(OnlineUsers, userName)
	//fmt.Println("OnlineUsers",OnlineUsers)
	chatlog.Std.WithFields(log.Fields{
		"username": userName,
	}).Info("Append The user to OnLine_user")
}

//DeleteUser 从OnlineUser中删除一个在线用户
func DeleteUser(userName string) {
	/*
		方法一：
		用户下线断开的是通讯地址，即 net.conn，因此维护在线用户列表和在线用户信息时
		需要根据这个net.conn 进行反查获取到用户名，进而进行信息删除。
		方法二：
		直接根据用户名删除，每个goroutine 都包含有userName，理论上不会删错信息。
	*/
	fmt.Println("online delete username ", userName)
	//删除 OnlineUserInfo map 中的信息
	delete(OnlineUserInfo, userName)
	//从OnlineUsers删除下线的用户
	if len(OnlineUsers) <= 1 {
		OnlineUsers = OnlineUsers[0:0]
	} else {
		for i := 0; i < len(OnlineUsers); i++ {
			if OnlineUsers[i] == userName {
				OnlineUsers = append(OnlineUsers[:i], OnlineUsers[i+1:]...)
			}
		}
	}
	fmt.Printf("delete OnlineUsers: \t%v\n OnlineUserInfo: \t%v\n ", OnlineUsers, OnlineUserInfo)
}

// SetNotOnLine 添加用户到未在线、

func SetNotOnLine(userName string) {

	NotOnline = append(NotOnline, userName)
}

//GetUserStatus 获取当前用户是否在线，如果存在则返回TRUE
/*
	userList 全局在线用户列表
	userName 需要查询的用户
*/
func GetUserStatus(userName string, userList []string) (b bool) {
	for _, value := range userList {
		if userName == value {
			b = true
			break
		} else {
			b = false
			continue
		}
	}
	return b
}
