package main

import (
	"chat/client/begin"
	"chat/client/worker"
)

func init() {
	begin.Connecting = begin.C_connting()

}

func main() {
	wk := &worker.Work{
		Conn: begin.Connecting,
	}
	wk.Worker()
}
