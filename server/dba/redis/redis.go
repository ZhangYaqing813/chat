package redism

import (
	msg "chat/Message_type"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var MyRedis *RedisOpt

var RedisPool *redis.Pool

type RedisOpt struct {
	pool *redis.Pool
}

func RedisPools(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	RedisPool = &redis.Pool{

		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func RedisFac(pool *redis.Pool) (MyRedis *RedisOpt) {
	MyRedis = &RedisOpt{
		pool: pool,
	}
	fmt.Println("redis init running", pool)
	return
}

// AddMessage 实现信息缓存到redis
/*
	1、redis 缓存的key 的命名方式为 key:userA field: recv_from_user/send_to_user/unread_message value
		reFromUser 		该user 接收到的信息。
		sendToUser   	该user 发送的信息
		unreadMessage 	该user 未读的信息
	2、确定 存入用户的那个一个field 中
*/
func (R *RedisOpt) AddMessage(username string, messages msg.Messages, fieldName string) (code int, err error) {
	//申请连接
	conn := R.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("redis conn close failed ", err)
		}
	}()
	// redis 认证，认证通过后进行后续操作
	_, err = conn.Do("AUTH", "12345678")
	if err != nil {
		//fmt.Println("Redis auth failed ", err)
		return code, err
	} else {
		//写入redis
		_, err = conn.Do("HSET", username, fieldName, messages)
		//fmt.Println("HSET______",conn_user)
		if err != nil {
			code = msg.FAILED
			return code, err
		}
		code = msg.SUCCESS
	}
	return code, err
}

//GetMessage

func (R *RedisOpt) GetMessage(userName string, fieldName []string) (messages []string, err error) {
	fmt.Println("redis get user ", userName)
	//申请连接
	index := 0

	conn := R.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("redis conn close failed ", err)
		}
	}()

	_, err = conn.Do("AUTH", "12345678")
	if err != nil {
		return
	} else {
		for _, key := range fieldName {
			cAdd, err := redis.String(conn.Do("HGET", userName, key))
			if err != nil {
				continue
			} else {
				//onlineuser[index].UserName=key
				//onlineuser[index].UserConn = cAdd
				fmt.Println("redis get user cAdd", cAdd)
				messages = append(messages, cAdd)
			}
			index++
		}
	}

	return messages, err
}

func (R *RedisOpt) Delete(username string, fieldNmae string) (code int, err error) {
	//申请连接
	conn := R.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("redis conn close failed ", err)
		}
	}()

	_, err = conn.Do("AUTH", "12345678")
	if err != nil {
		fmt.Println("Redis auth failed ", err)
		return msg.FAILED, err
	} else {
		for _, key := range username {
			_, err := conn.Do("HDEL", "OnLine_user", key)
			if err != nil {
				continue
			}
		}
	}
	return msg.SUCCESS, err
}
