package gateway

import (
	msg "chat/Message_type"
	mySql "chat/server/dba/mysql"
	"chat/server/msgproc"
	"chat/server/untils"
	_ "encoding/hex"
	"encoding/json"
	"fmt"
	_ "net"
)

type Gateway struct {
	msgproc.Messager
}

//

func (G *Gateway) Gateway() {
	slr := &untils.Slr{
		MysqlConnect: mySql.MysqlConnect{
			DB: mySql.MSconn,
		},
	}
	//var message msg.Messages
	//	go func() {
	for {
		//fmt.Println("go func()",G.Conn)
		message, err := G.MsgReader()
		if err != nil {
			fmt.Println("err")
			break
		}

		switch message.Type {
		// 处理用户登录逻辑

		case msg.LoginMsgType:
			var userinfo msg.LoginMsg
			var resPoneLoginMsg msg.LResMsg
			var lMsg msg.Messages
			// 1、解析message 中的 message.data 字段
			err = json.Unmarshal([]byte(message.Data), &userinfo)
			if err != nil {
				fmt.Println(err)
			}
			// 2、还原后传参给SLogin  处理
			// 用户登录验证成功后返回给客户一个登录成功的code
			// 单独使用一个code 变量，用于后续有需要的考虑
			code := slr.Slogin(userinfo)
			// 3、接收sLogin 返回数据组装返回给client
			resPoneLoginMsg.Code = code
			lMsg.Type = msg.RegMsgType
			lMsg.Data = string(G.Msgjson(resPoneLoginMsg))
			// 4、回复客户端认证结果
			G.MsgSender(lMsg)
			//后续操作
			if code == msg.SUCCESS {
				//通知其他用户改用户上线成功
				// 当OnlineUsers 长度为零是表示没有在线用户，不需要发送用户上线通知
				if untils.OnlineUsers[0] == "" {
					//上线用户加入在线用户列表
					untils.AddUser(userinfo.UserName, G.Messager.Conn)
				} else {
					// 通知在线用户
					G.NotifyOnline(userinfo.UserName, true)
				}
			}
			continue
		// 处理用户注册逻辑
		// 1、解析message 中的 message.data 字段
		// 2、还原后传参给Register 处理
		// 3、接收register 返回数据组装返回给client
		case msg.ResMsg:
			// register func
			var userinfo msg.RegMsg
			var resPoneLoginMsg msg.LResMsg
			var rRMsg msg.Messages
			err := json.Unmarshal([]byte(message.Data), &userinfo)
			if err != nil {
				fmt.Println(err)
			}
			code, err := slr.Register(userinfo)
			if err != nil {
				fmt.Println("register failed ", err)
			}
			// 其他处理逻辑
			//........
			resPoneLoginMsg.Code = code
			rRMsg.Type = msg.RegMsgType
			rRMsg.Data = string(G.Msgjson(resPoneLoginMsg))
			G.MsgSender(rRMsg)

		// 用户消息转发处理逻辑
		//
		case msg.ChatMode:
			var dia msg.Dialogue
			err = json.Unmarshal([]byte(message.Data), &dia)
			if err != nil {
				fmt.Println(err)
			}
			G.Transmit(dia, message)
		default:
			fmt.Println("and so on ......")
		}
	}
	//	}()

}
