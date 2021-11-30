package untils

import (
	msg "chat/Message_type"
	"fmt"
)

type Slr struct {
}

func (S *Slr) Slogin(userinfo msg.LoginMsg) (code int) {
	userid := 100
	userpwd := "abc"

	if userinfo.UserID == userid && userinfo.UserPwd == userpwd {
		fmt.Println("登录成功")
		code = msg.SUCCESS
	} else {
		code = msg.FAILED
	}

	return code
}
