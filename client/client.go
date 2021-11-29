package main

import (
	"chat/client/init"
	"chat/client/worker"
)

func init() {
	init.Connecting = init.C_connting()

}

func main() {
	wk := &worker.Work{
		Conn: init.Connecting,
	}
	wk.Worker()
}
