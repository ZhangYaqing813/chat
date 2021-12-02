package memnu

import (
	messagetype "chat/Message_type"
	"chat/client/client_func"
	"fmt"
	"os"
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
			}

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
