package main

import (
	"chat/Pb_mothd/msgproc"
	msconnecting "chat/server/dba/mysql"
	"chat/server/gateway"
	"fmt"
	"net"
)

func init() {
	//初始化工作
	msconnecting.MSconn = msconnecting.Factroy()
}

func main() {
	//var b [8192]byte
	lister, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("net.listen failed ", err)
	}
	defer lister.Close()
	for {
		conn, err := lister.Accept()
		if err != nil {
			fmt.Println("linster.accept failed ", err)
		}

		go func() {
			gw := &gateway.Gateway{
				Messager: msgproc.Messager{
					Conn: conn,
				},
			}
			gw.Gateway()
		}()
	}
}
