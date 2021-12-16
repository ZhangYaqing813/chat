package main

import (
	chatlog "chat/chatLog"
	"chat/client/begin"
	"chat/client/worker"
)

func init() {
	//初始化 client 到 server 的通讯
	begin.Connecting = begin.C_connting()
	chatlog.Init()

}

func main() {
	wk := &worker.Work{
		Conn: begin.Connecting,
	}
	wk.Worker()
}
