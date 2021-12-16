package msconnecting

import (
	message_type "chat/Message_type"
	chatlog "chat/chatLog"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"os"
)

var MSconn *sqlx.DB

type MysqlConnect struct {
	DB *sqlx.DB
}

func Factroy() *sqlx.DB {
	// 修改为参数，或配置文件
	dbconnect, err := sqlx.Open("mysql", "root:123456@tcp(172.30.1.251:3306)/chat")
	if err != nil {
		chatlog.Std.Fatalf("MySql connect failed %v", err)
		os.Exit(1)
	}
	chatlog.Std.Info("MySql 初始化完成")
	return dbconnect
}

// Insert  插入数据
func (M *MysqlConnect) Insert(msg message_type.RegMsg) (id int64, err error) {

	r, err := M.DB.Exec("insert into chat.users(username,password)values (?,?)", msg.UserName, msg.Password)
	if err != nil {
		chatlog.Std.WithFields(log.Fields{
			"username": msg.UserName,
			"password": "******",
		}).Fatalf("MySql Insert  failed %v", err)
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
	err := M.DB.Select(&userinfo, "select userid,username, password from chat.users where username =? ", msg.UserName)
	if err != nil {
		chatlog.Std.WithFields(log.Fields{
			"username": msg.UserName,
			"password": "******",
		}).Fatalf("MySql search failed %v", err)
		return
	}
	return userinfo
}
