package gateway

import (
	msg "chat/Message_type"
	"chat/Pb_mothd/msgproc"
	msconnecting "chat/server/dba/mysql"
	"chat/server/untils"
	"encoding/json"
	"fmt"
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
	var message msg.Messages

	message = G.MsgReader()
	fmt.Println("message = G.MsgReader()", message)
	switch message.Type {

	case msg.LoginMsgType:
		//login func
		var userinfo msg.LoginMsg
		var responeloginmsg msg.LResMsg
		var lmsg msg.Messages
		json.Unmarshal([]byte(message.Data), &userinfo)
		fmt.Println("msg.LoginMsgType userinfo", userinfo)
		code := slr.Slogin(userinfo)

		responeloginmsg.Code = code
		lmsg.Type = msg.RegMsgType
		lmsg.Data = string(G.Msgjson(responeloginmsg))
		G.MsgSender(lmsg)

	case msg.ResMsg:
		// register func
		var userinfo msg.RegMsg
		var responeregmsg msg.LResMsg
		var rrmsg msg.Messages
		code, err := slr.Register(userinfo)
		if err != nil {
			fmt.Println("register failed ", err)
		}

		responeregmsg.Code = code
		rrmsg.Type = msg.RegMsgType
		rrmsg.Data = string(G.Msgjson(responeregmsg))
		G.MsgSender(rrmsg)

	default:
		fmt.Println("and so on ......")

	}

}
