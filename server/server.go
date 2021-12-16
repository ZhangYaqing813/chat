package main

import (
	chatlog "chat/chatLog"
	msconnecting "chat/server/dba/mysql"
	redism "chat/server/dba/redis"
	"chat/server/gateway"
	"chat/server/msgproc"
	log "github.com/sirupsen/logrus"
	"net"
)

/*
server.go 是整个服务端的入口，运行后这个服务端就运行起来了。
server.go 主要包括两个部分，一部分是基础环境的初始化，比如连接redis,mysql,以及读取配置文件等（配置文件目前不支持）
第二部分是main 函数，主要功能是就是建立监听socket，调用相关方法处理client 发送的数据
*/

func init() {
	//初始化工作
	// 初始化mysql 链接
	msconnecting.MSconn = msconnecting.Factroy()
	//初始话redis
	redism.RedisPools("172.30.1.2:6379", 16, 0, 300)
	redism.MyRedis = redism.RedisFac(redism.RedisPool)
	chatlog.Init()
}

func main() {

	//打开监听地址
	lister, err := net.Listen("tcp", "127.0.0.1:19000")
	if err != nil {
		chatlog.Std.Fatal(err)
	}

	//接收客户端请求
	for {
		conn, err := lister.Accept()

		if err != nil {
			chatlog.Std.Fatal(err)
		}
		chatlog.Std.WithFields(log.Fields{
			"RemoteIP": conn.RemoteAddr(),
		}).Info("Client connected ")
		// 多协程处理客户端请求
		go func() {
			//初始化路由实例，并将conn 地址传递
			gw := &gateway.Gateway{
				Messager: msgproc.Messager{
					Conn: conn,
				},
			}
			//调用gateway 方法
			gw.Gateway()
		}()
	}
}
