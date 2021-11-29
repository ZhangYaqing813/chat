package worker

import (
	"net"
)

type Work struct {
	Conn net.Conn
}

// Worker 主要用于和server 端的数据通信。
func (W *Work) Worker() {
	//msgp := &msgproc.Messager{
	//	W.Conn,
	//}

}
