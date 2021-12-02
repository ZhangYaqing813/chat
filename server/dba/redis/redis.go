package redism

import (
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
	return
}

func (R *RedisOpt) Add(db, data string, username string) (code int, err error) {
	//申请连接
	conn := R.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("redis conn close failed ", err)
		}
	}()

	//_, err := conn.Do("HSET","")

	return
}
