package gateway

import (
	msg "chat/Message_type"
	msconnecting "chat/server/dba/mysql"
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
		msconnecting.MysqlConnect{
			DB: msconnecting.MSconn,
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
		//fmt.Println("gateway recv msg ", message)
		switch message.Type {
		// 处理用户登录逻辑
		// 1、解析message 中的 message.data 字段
		// 2、还原后传参给Slogin  处理
		// 3、接收Slogin 返回数据组装返回给client
		// 4、根据返回的code 	进行下一步处理
		case msg.LoginMsgType:
			var userinfo msg.LoginMsg
			var responeloginmsg msg.LResMsg
			var lmsg msg.Messages
			json.Unmarshal([]byte(message.Data), &userinfo)
			// 用户登录验证成功后返回给客户一个登录成功的code
			code := slr.Slogin(userinfo)
			responeloginmsg.Code = code
			lmsg.Type = msg.RegMsgType
			lmsg.Data = string(G.Msgjson(responeloginmsg))
			G.MsgSender(lmsg)

			//通知其他用户改用户上线成功
			// 当untils.OnlineUsers 长度为零是表示没有在线用户，不需要发送用户上线通知
			if untils.OnlineUsers[0] == "" {
				//上线用户加入在线用户列表
				untils.AddUser(userinfo.UserName, G.Messager.Conn)
			} else {
				// 通知在线用户
				G.NotifyOnline(userinfo.UserName, true)
			}

			continue
		// 处理用户注册逻辑
		// 1、解析message 中的 message.data 字段
		// 2、还原后传参给Register 处理
		// 3、接收register 返回数据组装返回给client
		case msg.ResMsg:
			// register func
			var userinfo msg.RegMsg
			var responeregmsg msg.LResMsg
			var rrmsg msg.Messages
			json.Unmarshal([]byte(message.Data), &userinfo)
			code, err := slr.Register(userinfo)
			if err != nil {
				fmt.Println("register failed ", err)
			}
			// 其他处理逻辑
			//........
			responeregmsg.Code = code
			rrmsg.Type = msg.RegMsgType
			rrmsg.Data = string(G.Msgjson(responeregmsg))
			G.MsgSender(rrmsg)

		// 用户消息转发处理逻辑
		//
		case msg.ChatMode:
			var dia msg.Dialogue
			json.Unmarshal([]byte(message.Data), &dia)
			G.Transmit(dia, message)

		default:
			fmt.Println("and so on ......")
		}
	}
	//	}()

}
