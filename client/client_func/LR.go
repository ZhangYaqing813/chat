package client_func

import (
	messagetype "chat/Message_type"
	"chat/Pb_mothd/msgproc"
	"fmt"
)

type LR struct {
	//conn net.Conn
	msgproc.Messager
}

func (L *LR) Login(loginmsg messagetype.LoginMsg) {
	var msg messagetype.Messages

	msg.Type = messagetype.LoginMsgType
	msg.Data = string(L.Msgjson(loginmsg))

	L.MsgSender(msg)

	msg = L.MsgReader()
	fmt.Println(msg)

}
