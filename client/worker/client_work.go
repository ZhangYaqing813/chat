package worker

import (
	"chat/client/client_func"
	"chat/client/memnu"
	"chat/client/msgproc"
	"net"
)

type Work struct {
	Conn net.Conn
	memnu.Menus
}

// Worker 主要用于和server 端的数据通信。
func (W *Work) Worker() {

	wk := &memnu.Menus{
		client_func.LR{
			msgproc.Messager{
				W.Conn,
			},
		},
	}
	wk.Firstlevel()
}
