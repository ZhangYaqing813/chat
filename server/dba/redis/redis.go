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

func (R *RedisOpt) Add(onlineUser string, username string) (code int, err error) {
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
		return
	} else {

		_, err = conn.Do("HSET", "OnLine_user", username, onlineUser)
		//fmt.Println("HSET______",conn_user)
		if err != nil {
			fmt.Println("add online user failed ", err)
			code = msg.FAILED
			return code, err

		}
		code = msg.SUCCESS
	}
	return code, err
}

func (R *RedisOpt) Get(username []string) (user []string, err error) {
	//申请连接
	index := 0
	onlineuser := make([]string, 1024)
	conn := R.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("redis conn close failed ", err)
		}
	}()

	_, err = conn.Do("AUTH", "12345678")
	if err != nil {
		return user, err
	} else {

		for _, key := range username {
			cAdd, err := redis.String(conn.Do("HGET", "OnLine_user", key))
			if err != nil {
				continue
			} else {
				//onlineuser[index].UserName=key
				//onlineuser[index].UserConn = cAdd
				user = append(user, cAdd)
			}
			index++
		}
	}

	return onlineuser, err
}

func (R *RedisOpt) Delete(username []string) (code int, err error) {
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
