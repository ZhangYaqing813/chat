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

/*
	redis 缓存的key 的命名方式为 key:userAReFromUserB/userASendToUserB/userAUnreadMessage field:时间戳： value：message
    key:
        userAReFromUserB 	该user 接收到的信息。
        userASendToUser B  	该user 发送的信息
        userAUnreadMessage 	该user 未读的信息
*/

// AddMessage 实现信息缓存到redis

func (R *RedisOpt) AddMessage(key string, messages msg.Messages, fieldName string) (code int, err error) {
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
		_, err = conn.Do("HSET", key, fieldName, messages)
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

func (R *RedisOpt) GetMessage(key string, fieldName []string) (messages []string, err error) {
	fmt.Println("redis get user ", key)
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
		for _, field := range fieldName {
			//改进一下，获取多个field

			value, err := redis.String(conn.Do("HGET", key, field))
			if err != nil {
				continue
			} else {
				fmt.Println("redis get user cAdd", value)
				messages = append(messages, value)
			}
			index++
		}
	}

	return messages, err
}

func (R *RedisOpt) Delete(key string) (code int, err error) {
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
		//for _, valud := range username {
		//	_, err := conn.Do("HDEL", key, key)
		//	if err != nil {
		//		continue
		//	}
		//}
		// 删除 userAUnreadMessage 这里面的所有信息。
		// _, err = conn.Do("Hdel", key,)

	}
	return msg.SUCCESS, err
}
