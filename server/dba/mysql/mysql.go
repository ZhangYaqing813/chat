package msconnecting

import (
	message_type "chat/Message_type"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MSconn *sqlx.DB

type MysqlConnect struct {
	DB *sqlx.DB
}

func Factroy() *sqlx.DB {
	dbconnect, err := sqlx.Open("mysql", "root:root@tcp(172.30.1.2:3306)/chat")

	if err != nil {
		fmt.Println("mysql db connect failed", err)
	}
	fmt.Println("mysql factory running")
	//defer dbconnect.Close()

	return dbconnect
}

// Insert  插入数据
func (M *MysqlConnect) Insert(msg message_type.RegMsg) (id int64, err error) {

	r, err := M.DB.Exec("insert into chat.users(username,password)values (?,?)", msg.UserName, msg.UserPwd)
	if err != nil {
		fmt.Println("insert data failed ", err)
		return 0, err
	}

	userid, err := r.LastInsertId()
	if err != nil {
		fmt.Println("返回userid failed ", err)
		return 0, err
	}
	return userid, err
}

// Update 更新数据
func (M *MysqlConnect) Update() {

}

// Delete 删除数据
func (M *MysqlConnect) Delete() {

}

// Select 查找数据

func (M *MysqlConnect) Select(msg message_type.LoginMsg) (userinfo []message_type.LoginMsg) {

	fmt.Println("msg message_type.LoginMsg Select ", msg)
	err := M.DB.Select(&userinfo, "select userid,username, password from chat.users where username =? ", msg.UserName)
	if err != nil {
		fmt.Println("exec failed ", err)
		return
	}
	fmt.Println(userinfo)
	return userinfo
}
