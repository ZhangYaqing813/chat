package memnu

import (
	msg "chat/Message_type"
	"chat/client/client_func"
	"fmt"
	"os"
)

var userinfo msg.LoginMsg

type Menus struct {
	client_func.LR
}

func (M *Menus) Firstlevel() {
	//var message msg.Messages
	var loginmsg msg.LoginMsg
	var key int

	for {
		fmt.Println("\t\t******** 欢迎来到聊天室 ********")
		fmt.Println("\t\t********1、用户登录  ********")

		fmt.Println("\t\t********2、用户注册  ********")
		fmt.Println("\t\t********3、退出系统  ********")
		fmt.Println("请输入选择：(1-3):")
		fmt.Scanf("%d\n", &key)

		//user := 100
		//password := "zyq"
		switch key {
		case 1:
			//var userinfo PublicMethods.LoginMsg
			fmt.Println("请输入用户ID：")
			fmt.Scanf("%d\n", &loginmsg.UserID)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &loginmsg.UserPwd)

			M.Login(loginmsg)

			//if loginmsg.UserID == user && loginmsg.UserPwd == password {
			//	fmt.Println("登录成功")
			//} else {
			//	fmt.Println("登录失败")
			//	return
			//}

			//err := usermsg.Login(userinfo)
			//if err != nil {
			//	fmt.Printf("user login failed %v", err)
			//	os.Exit(0)
			//}

		//case 2:
		//	//var usReg PublicMethods.RegMsg
		//	fmt.Println("新用户注册")
		//	fmt.Println("请输入用户ID：")
		//	fmt.Scanf("%d\n", &usReg.UserID)
		//	fmt.Println("请输入用户密码：")
		//	fmt.Scanf("%s\n", &usReg.UserPwd)
		//
		//	err := usermsg.Register(usReg)
		//	if err != nil {
		//		fmt.Println("user register failed ", err)
		//	}
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
