package untils

import (
	msg "chat/Message_type"
	chatlog "chat/chatLog"
	msconnecting "chat/server/dba/mysql"
	log "github.com/sirupsen/logrus"
)

type Slr struct {
	msconnecting.MysqlConnect
}

// Slogin 用户首次登录暂时需要提供用户名和密码，后面可以根据实际需求修改为用户名、用户ID以及用户邮箱的方式
// 扩展：用户登录成功时可以生成一个session 表项，用户ID，client端与server端建立连接的内存地址，并且存入Redis，以备后用
// 用户登录成功后，可以进行Redis 的缓存，并设置相应的过期时间，减少sql 查询

func (S *Slr) Slogin(userinfo msg.LoginMsg) (code int) {
	// 查询数据库
	user := S.Select(userinfo)
	// 判断用户信息是否匹配
	if user[0].UserName == userinfo.UserName && user[0].Password == userinfo.Password {
		//记录日志
		chatlog.Std.WithFields(log.Fields{
			"username": userinfo.UserName,
			"userid":   userinfo.UserID,
		}).Info("上线")
		//设置code 值
		code = msg.SUCCESS
	} else {
		// 如果不匹配则返回相应的值
		code = msg.FAILED
	}
	return code
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
