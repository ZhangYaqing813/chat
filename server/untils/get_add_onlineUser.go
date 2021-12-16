package untils

import (
	chatlog "chat/chatLog"
	log "github.com/sirupsen/logrus"
	"net"
)

// 维护一个在线用的map，用户获取在线用户

var OnlineUserInfo = make(map[string]net.Conn)
var OnlineUsers = make([]string, 1024)

// AddUser 向 OnlineUser 中添加一个上线的用户
func AddUser(userName string, conn net.Conn) {

	chatlog.Std.WithFields(log.Fields{
		"username": userName,
	}).Info("Append The user to OnLine_user")
	OnlineUserInfo[userName] = conn
	for user, _ := range OnlineUserInfo {
		OnlineUsers = append(OnlineUsers, user)
	}
}

//DeleteUser 从OnlineUser中删除一个在线用户
func DeleteUser(userName string) {
	delete(OnlineUserInfo, userName)
	if len(OnlineUsers) <= 1 {
		OnlineUsers = OnlineUsers[0:0]
	} else {
		for i := 0; i < len(OnlineUsers); i++ {
			if OnlineUsers[i] == userName {
				OnlineUsers = append(OnlineUsers[:i], OnlineUsers[i+1:]...)
			}
		}
	}

}
