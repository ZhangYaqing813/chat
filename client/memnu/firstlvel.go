package memnu

import (
	messagetype "chat/Message_type"
	"chat/client/client_func"
	"fmt"
	"os"
	"strings"
	"time"
)

//var userinfo msg.LoginMsg

type Menus struct {
	client_func.LR
}

func (M *Menus) Firstlevel() {
	//var message msg.Messages
	var loginmsg messagetype.LoginMsg
	var userReg messagetype.RegMsg
	var key int

	for {
		fmt.Println("\t\t******** 欢迎来到聊天室 ********")
		fmt.Println("\t\t********1、用户登录  ********")

		fmt.Println("\t\t********2、用户注册  ********")
		fmt.Println("\t\t********3、退出系统  ********")
		fmt.Println("请输入选择：(1-3):")
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			//var userinfo PublicMethods.LoginMsg
			fmt.Println("请输入用户名：")
			fmt.Scanf("%s\n", &loginmsg.UserName)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &loginmsg.Password)

			code, error := M.Login(loginmsg)
			if code == 200 && error == "" {
				fmt.Println("用户登录成功 ")
				M.loginLevel(loginmsg.UserName)
			}
			// second level

		case 2:
			fmt.Println("新用户注册")
			fmt.Println("请输入用户名：")
			fmt.Scanf("%s\n", &userReg.UserName)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userReg.Password)
			code, error := M.Register(userReg)
			if code == 200 && error == "" {
				fmt.Println("注册成功")
			}

		//case 3:
		//	err := usermsg.ExitOS(userinfo.UserID)
		//	if err != nil {
		//		fmt.Printf("exit os failed err:%v\n", err)
		//	}
		default:
			fmt.Println("暂不处理")
			os.Exit(0)
		}

	}

}

func (M *Menus) loginLevel(loginUser string) {
	var skey int
	var dialogue messagetype.Dialogue
	var name string
	fmt.Println("恭喜登录成功,请选择需要的功能")
	fmt.Println("1、单聊")
	fmt.Println("2、群聊")
	fmt.Println("3、更新在线用户")
	fmt.Println("4、退出")
	fmt.Println("请输入选项：")
	fmt.Scanf("%d\n", &skey)
	switch skey {

	case 1:
		// 配置 聊天模式位单聊，
		dialogue.SendMod = messagetype.SINGLE
		// 消息发送对象
		fmt.Println("请输入接收人")
		fmt.Scanf("%s\n", &name)
		dialogue.ToUsers = append(dialogue.ToUsers, name)
		dialogue.Sender = loginUser

		// 输入内用，

		fmt.Println("请输入内容")
		fmt.Scanf("%s\n", &dialogue.Content)
		dialogue.SendTime = time.Now().Format("2006-01-02 15:04:05")
		// 封装 聊天信息
		M.Chat(dialogue)
		fmt.Println(dialogue)

	case 2:
		// 配置 发送模式为多个，
		dialogue.SendMod = messagetype.MULTIPLE
		// 消息发送对象
		fmt.Println("请输入接收人,以逗号分开")
		fmt.Scanf("%s\n", &name)
		sep := ","
		dialogue.ToUsers = strings.Split(name, sep)
		fmt.Println(dialogue.ToUsers)
		dialogue.Sender = loginUser

		// 输入内用，
		fmt.Println("请输入内容")
		fmt.Scanf("%s\n", &dialogue.Content)
		dialogue.SendTime = time.Now().Format("2006-01-02 15:04:05")
		// 封装 聊天信息
		fmt.Println(dialogue)
	case 3:
		// 更新在线用户
		//message.Type = messagetype.UPDATE
		//fmt.Println(message)

	case 4:
		//用户退出

	}

}
