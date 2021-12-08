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
		for {
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
				//fmt.Println("msg.LoginMsgType userinfo", userinfo)
				code := slr.Slogin(userinfo)
				responeloginmsg.Code = code

				lmsg.Type = msg.RegMsgType
				lmsg.Data = string(G.Msgjson(responeloginmsg))
				G.MsgSender(lmsg)
				//fmt.Println("&G.Conn",&G.Conn)
				//其他处理逻辑

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
				continue
			case msg.ChatMode:
				var dia msg.Dialogue
				json.Unmarshal([]byte(message.Data), &dia)

				G.Transmit(dia, message)

			default:
				fmt.Println("and so on ......")

			}
		}

	}()

}
