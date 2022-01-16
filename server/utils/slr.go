package utils

import (
	msg "chat/Message_type"
	chatlog "chat/chatLog"
	msconnecting "chat/server/dba/mysql"
	redism "chat/server/dba/redis"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Slr struct {
	msconnecting.MysqlConnect
}

// Slogin 用户首次登录暂时需要提供用户名和密码，后面可以根据实际需求修改为用户名、用户ID以及用户邮箱的方式
// 扩展：用户登录成功时可以生成一个session 表项，用户ID，client端与server端建立连接的内存地址
// 用户登录成功后，可以进行Redis 的缓存，并设置相应的过期时间，减少sql 查询

func (S *Slr) Slogin(userinfo msg.LoginMsg) (code int, unMessages []string) {

	// 查询数据库

	user, _ := S.Select(userinfo.UserName)
	// 判断用户信息是否匹配

	if user[0].UserName == userinfo.UserName && user[0].Password == userinfo.Password {
		//记录日志
		chatlog.Std.WithFields(log.Fields{
			"username": userinfo.UserName,
			"userid":   userinfo.UserID,
		}).Info("上线")
		//设置code 值
		code = msg.SUCCESS
		// 增加一个步骤，当用户认证成功后，增加一个发送未读消息的功能。

	} else {
		// 如果不匹配则返回相应的值
		code = msg.FAILED
	}
	// 用户登录后去查询redis 是否有未读信息，如果有则发送。
	if GetUserStatus(userinfo.UserName, NotOnline) {
		unMessages, err := redism.MyRedis.GetMessage(userinfo.UserName)
		if err == nil {
			return code, unMessages
		}
	}

	return
}

// Register  用户注册处理功能
func (S *Slr) Register(regMsg msg.RegMsg) (code int, err error) {

	id, err := S.Insert(regMsg)
	if err != nil || id == 0 {
		chatlog.Std.WithFields(log.Fields{
			"username": regMsg.UserName,
		}).Fatal("注册失败")
		return
	} else {
		code = msg.SUCCESS
	}
	return code, err
}

// Modify 更新用户信息（主要针对更新用户信息，密码或者email）
func (S *Slr) Modify(modify msg.UserUpdate, userName string) (code int, err error) {
	// 检查当前用户是否在线，如果不在线怎不能进行信息的修改
	if GetUserStatus(userName, OnlineUsers) {
		res, err := S.Update(modify, userName)
		if res > 0 {
			chatlog.Std.Info("modify user info OK")
			code = msg.SUCCESS
		} else {
			code = msg.FAILED
			chatlog.Std.Error("modify user info failed ", err)
		}
	} else {
		code = msg.FAILED
		err = fmt.Errorf("当前用户未登录，请登录后进行修改")
	}
	return code, err
}
