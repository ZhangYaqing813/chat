package gateway

import (
	msg "chat/Message_type"
	"chat/Pb_mothd/msgproc"
	"chat/server/untils"
	"encoding/json"
	"fmt"
)

type Gateway struct {
	msgproc.Messager
	untils.Slr
}

//

func (G *Gateway) Gateway() {
	var message msg.Messages

	message = G.MsgReader()

	switch message.Type {

	case msg.LoginMsgType:
		//login func
		var userinfo msg.LoginMsg
		var responeloginmsg msg.LResMsg
		var lmsg msg.Messages
		json.Unmarshal([]byte(message.Data), &userinfo)
		code := G.Slogin(userinfo)

		responeloginmsg.Code = code
		lmsg.Type = msg.RegMsgType
		lmsg.Data = string(G.Msgjson(responeloginmsg))
		G.MsgSender(lmsg)

	case msg.ResMsg:
	// register func
	//:

	default:
		fmt.Println("and so on ......")

	}

}
