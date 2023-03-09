package db

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool // 创建redis连接池

func init() {
	pool = &redis.Pool{ //实例化一个连接池
		MaxIdle:   10,      //最初的连接数量
		MaxActive: 1000000, //最大连接数量
		// MaxActive:   0,                 //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: time.Second * 240, //连接关闭时间 480秒 （480秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			conn, err := redis.Dial("tcp", "192.168.195.174:9222")
			if err != nil {
				return nil, err
			}
			// if _, err := conn.Do("AUTH", "123456"); err != nil {
			// if _, err := conn.Do("AUTH", ""); err != nil {
			// 	conn.Close()
			// 	return nil, err
			// }
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			// if time.Since(t) < time.Minute {
			// 	return nil
			// }
			_, err := conn.Do("PING")
			return err
		},
	}
}

func GetRedisDbInstance() redis.Conn {
	return pool.Get()
}
