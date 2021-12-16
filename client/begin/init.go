package begin

import (
	"fmt"
	"net"
)

var Connecting net.Conn

func C_connting() (conn net.Conn) {
	conn, err := net.Dial("tcp", "127.0.0.1:19000")
	if err != nil {
		fmt.Println("client connect failed ", err)
		return
	}
	return conn
}
