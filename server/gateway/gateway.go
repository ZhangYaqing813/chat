package gateway

import (
	msg "chat/Message_type"
	mySql "chat/server/dba/mysql"
	"chat/server/msgproc"
	"chat/server/untils"
	_ "encoding/hex"
	"encoding/json"
	"fmt"
	_ "net"
)

type Gateway struct {
	msgproc.Messager
}

func (G *Gateway) Gateway() {
	slr := &untils.Slr{
		MysqlConnect: mySql.MysqlConnect{
			DB: mySql.MSconn,
		},
	}
	var username_tmp string // 用于删除在线用户
	//var message msg.Messages
	//	go func() {
	for {
		message, err := G.MsgReader()
		// 增加一步对客户端的处理，用于处理下线或直接断开的客户端，然后通知在线客户某用户下线
		/*
			客户端下线后 msgReader 会返回一个错误，利用这个错误来进行客户端下线的处理
			方法一：
				当net.conn.read 函数读取数据到字节流末尾时会返回一个 io.EOF 的问题，此时只是表示数据
				读取完毕，而不一定是client 断开，因此要做进一步处理，不在本次考虑范围内
			方法二：
				利用当前goroutine 中保存的userinfo.userName,这里就需要先进行 message 的解析然后保存userName

		*/
		if err != nil {
			fmt.Println(" G.MsgReader() err", err)

			G.NotifyOnline(username_tmp, false)
			untils.DeleteUser(username_tmp)
			break
		}
		switch message.Type {
		// 处理用户登录逻辑

		case msg.LoginMsgType:
			var userinfo msg.LoginMsg
			var resPoneLoginMsg msg.LResMsg
			var lMsg msg.Messages
			// 1、解析message 中的 message.data 字段
			err = json.Unmarshal([]byte(message.Data), &userinfo)
			if err != nil {
				fmt.Println(err)
			}
			username_tmp = userinfo.UserName
			// 2、还原后传参给SLogin  处理
			// 用户登录验证成功后返回给客户一个登录成功的code
			// 单独使用一个code 变量，用于后续有需要的考虑
			code := slr.Slogin(userinfo)
			// 3、接收sLogin 返回数据组装返回给client
			resPoneLoginMsg.Code = code
			lMsg.Type = msg.RegMsgType
			lMsg.Data = string(G.Msgjson(resPoneLoginMsg))
			// 4、回复客户端认证结果
			err = G.MsgSender(lMsg)
			if err != nil {
				fmt.Println("调用失败")
			}
			//后续操作
			if code == msg.SUCCESS {
				//通知其他用户改用户上线成功
				// 当OnlineUsers 长度为零是表示没有在线用户，不需要发送用户上线通知
				if untils.OnlineUsers[len(untils.OnlineUsers)-1] == "" {
					//上线用户加入在线用户列表
					//fmt.Println("111111")
					//fmt.Println("G.Messager.Conn",userinfo.UserName,G.Messager.Conn)
					untils.AddUser(userinfo.UserName, G.Messager.Conn)
				} else {
					// 通知在线用户
					//fmt.Println("22222222")
					fmt.Println("当前shangxian用户信息", untils.OnlineUserInfo)
					/*
						问题描述：
							用户上线顺序 A-B-C
							在执行NotifyOnline，后执行 AddUser 时，NotifyOnline 通知逻辑时获取在线用户的通讯地址，然后赋值给G.conn,
							这会导致 A收到两份C上线的通知，B无法收到C上线通知。原因是，在B上线后通知A是修改了B的通讯地址 M.conn 使其成为A
							A的通讯地址 M.Conn = untils.OnlineUserInfo[user] （user = A）

						解决方法一：
							修改执行顺序：
								由 NotifyOnline -> AddUser  变成为 AddUser —> NotifyOnline
								弊端：自己会收到自己上线的通知
						解决方法二：
								直接修改msgProc 的传参方式，即修改调用方式，每个方法都增加一个 conn net.conn 的参数，用于控制
								消息发送对象。
						本次处理方法采用第一种，第二种修改范围比较大

					*/
					untils.AddUser(userinfo.UserName, G.Messager.Conn)
					G.NotifyOnline(userinfo.UserName, true)

				}
			}
			continue
		// 处理用户注册逻辑
		// 1、解析message 中的 message.data 字段
		// 2、还原后传参给Register 处理
		// 3、接收register 返回数据组装返回给client
		case msg.ResMsg:
			// register func
			var userinfo msg.RegMsg
			var resPoneLoginMsg msg.LResMsg
			var rRMsg msg.Messages
			err := json.Unmarshal([]byte(message.Data), &userinfo)
			if err != nil {
				fmt.Println(err)
			}
			code, err := slr.Register(userinfo)
			if err != nil {
				fmt.Println("register failed ", err)
			}
			// 其他处理逻辑
			//........
			resPoneLoginMsg.Code = code
			rRMsg.Type = msg.RegMsgType
			rRMsg.Data = string(G.Msgjson(resPoneLoginMsg))
			err = G.MsgSender(rRMsg)
			if err != nil {
				fmt.Println("调用失败", msg.ResMsg)
			}

		// 用户消息转发处理逻辑
		//
		case msg.ChatMode:
			var dia msg.Dialogue
			err = json.Unmarshal([]byte(message.Data), &dia)
			if err != nil {
				fmt.Println(err)
			}
			G.Transmit(dia, message)
		default:
			fmt.Println("and so on ......")
		}
	}
	//	}()

}
