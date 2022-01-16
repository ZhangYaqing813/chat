package msconnecting

import (
	message_type "chat/Message_type"
	chatlog "chat/chatLog"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MSconn *sqlx.DB

type MysqlConnect struct {
	DB *sqlx.DB
}

func Factroy(user, password, address, dbname, port string) *sqlx.DB {
	//修改了连接方式
	fmt.Println("port= ", port)
	socket := user + ":" + password + "@tcp" + "(" + address + ":" + port + ")" + "/" + dbname
	fmt.Println("socket = ", socket)
	dbconnect, err := sqlx.Open("mysql", socket)

	if err != nil {
		fmt.Println("mysql db connect failed", err)
	}
	fmt.Println("mysql factory running")
	return dbconnect
}

// Insert  插入数据
func (M *MysqlConnect) Insert(msg message_type.RegMsg) (id int64, err error) {

	r, err := M.DB.Exec("insert into chat.users(username,password)values (?,?)", msg.UserName, msg.Password)
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
func (M *MysqlConnect) Update(modify message_type.UserUpdate, userName string) (result int64, err error) {
	// 根据modify 中的字段信息执行不同的更新语句
	var row int64
	chatlog.Std.Info("用户 修改密码", userName)
	switch modify.FieldName {
	case "password":
		res, err := M.DB.Exec("update chat.users set  password =? where username = ?", modify.NewContent, userName)
		if err != nil {
			chatlog.Std.Error("SQL Exec failed ", err)
		}
		row, _ = res.RowsAffected()

	case "email":
		res, err := M.DB.Exec("update chat.users set  email =? where username = ?", modify.NewContent, userName)
		if err != nil {
			chatlog.Std.Error("SQL Exec failed ", err)
		}
		row, _ = res.RowsAffected()

	default:
		fmt.Println("该字段不能被修改")
	}
	return row, err
}

// Delete 删除数据
func (M *MysqlConnect) Delete() {

}

// Select 查找数据

func (M *MysqlConnect) Select(userName string) (userinfo []message_type.LoginMsg, err error) {

	err = M.DB.Select(&userinfo, "select userid,username, password from chat.users where username =? ", userName)
	if err != nil {
		fmt.Println("exec failed ", err)
		//fmt.Println("111111111")
		return
	}
	//fmt.Println(userinfo)
	return userinfo, err
}
