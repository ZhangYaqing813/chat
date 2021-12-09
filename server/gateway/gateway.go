package gateway

import (
	msg "chat/Message_type"
	"chat/Pb_mothd/msgproc"
	msconnecting "chat/server/dba/mysql"
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
	//var onlineUser msg.UserOlineIntoRedis

	//var onlineUser  msg.UserOlineIntoRedis
	var message msg.Messages
	go func() {

		message = G.MsgReader()
		fmt.Println("message = G.MsgReader()", message)
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
			//其他处理逻辑
			//上线用户加入在线用户列表
			untils.AddUser(userinfo.UserName, G.Messager.Conn)
			//通知其他用户改用户上线成功

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

		case msg.ChatMode:
			var dia msg.Dialogue
			json.Unmarshal([]byte(message.Data), &dia)

			G.Transmit(dia, message)

		default:
			fmt.Println("and so on ......")

		}

	}()

}
