package client_func

import (
	message_type "chat/Message_type"
	"chat/Pb_mothd/msgproc"
	"fmt"
)

type LR struct {
	//conn net.Conn
	msgproc.Messager
}

func (L *LR) Login(loginmsg message_type.LoginMsg) {
	var msg message_type.Messages

	msg.Type = message_type.LoginMsgType
	msg.Data = string(L.Msgjson(loginmsg))

	L.MsgSender(msg)

	msg = L.MsgReader()
	fmt.Println(msg)

}
