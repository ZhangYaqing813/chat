package untils

import (
	msg "chat/Message_type"
	redism "chat/server/dba/redis"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
)

// 整体实现从redis 中取出onlineUser
// 用户登录成功后 将用户的信息缓存在redis，以便后续处理

//
func SetOnlineUsers(userName string, userConn net.Conn) (code int, err error) {
	// 将登录成功的用户内存地址存入到redis
	// 组装用户信息
	var onlineUser msg.UserOlineIntoRedis
	onlineUser.UserName = userName
	onlineUser.UserConn = userConn
	fmt.Println(onlineUser)
	// 序列化
	data, _ := json.Marshal(onlineUser)
	fmt.Println("data", data)
	// 存入redis 前进行一次编码，因为内存地址格式的特殊性，进行一次编码，
	//使用string()方法进行[]byte数组进行string 字符串进行转换后 onlineUser.UserConn 值为空
	//因此需要一次编码，本次只测试了 hex.EncodeToString 可以成功
	code, err = redism.MyRedis.Add(hex.EncodeToString(data), userName)
	return code, err
}

func GetOlineUser(users []string) (uOnline []msg.UserOlineIntoRedis) {
	var user msg.UserOlineIntoRedis
	tmp, _ := redism.MyRedis.Get(users)
	for index, _ := range tmp {
		buf, _ := hex.DecodeString(users[index-1])

		err := json.Unmarshal(buf, &user)
		if err != nil {
			fmt.Println("GetOlineUser Unmarshal failed ", err)
			return
		}
		uOnline = append(uOnline, user)
	}
	return uOnline
}
