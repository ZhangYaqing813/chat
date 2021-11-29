package main

import (
	"fmt"
	"net"
)

func init() {
	//初始化工作
}

func main() {

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
		fmt.Println(conn.RemoteAddr())
	}

}
