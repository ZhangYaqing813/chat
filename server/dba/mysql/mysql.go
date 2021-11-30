package msconnecting

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MSconn *sqlx.DB

type MysqlConnect struct {
	DB sqlx.DB
}

func Factroy() *sqlx.DB {
	dbconnect, err := sqlx.Open("mysql", "root:123456@tcp(172.30.1.2:3306)/chat")

	if err != nil {
		fmt.Println("mysql db connect failed", err)
	}
	fmt.Println("mysql factory running")
	dbconnect.Close()

	return dbconnect
}

// Insert  插入数据
func (M *MysqlConnect) Insert() {

}

// Update 更新数据
func (M *MysqlConnect) Update() {

}

// Delete 删除数据
func (M *MysqlConnect) Delete() {

}

// Modify 修改数据

func (M *MysqlConnect) Modify() {

}
