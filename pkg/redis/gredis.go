package gredis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

var RedisPool *redis.Pool

// Setup Initialize the Redis instance
func Setup() {
	/*
		设置 RedisConn 为 redis.Pool（连接池）并配置了它的一些参数：

		Dial：提供创建和配置应用程序连接的一个函数
		TestOnBorrow：可选的应用程序检查健康功能
		MaxIdle：最大空闲连接数
		MaxActive：在给定时间内，允许分配的最大连接数（当为零时，没有限制）
		IdleTimeout：在给定时间内将会保持空闲状态，若到达时间限制则关闭连接（当为零时，没有限制）
	*/
	//fmt.Println(setting.RedisSetting)
	var maxIdle = setting.RedisSetting.MaxIdle
	maxActive := setting.RedisSetting.MaxActive
	idleTimeout := setting.RedisSetting.IdleTimeout
	RedisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			/*c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err*/
			return redis.Dial("tcp", "192.168.10.33:6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

/*
文件内包含 Set、Exists、Get、Delete、LikeDeletes 用于支撑目前的业务逻辑，而在里面涉及到了如方法：
（1）RedisConn.Get()：在连接池中获取一个活跃连接
（2）conn.Do(commandName string, args ...interface{})：向 Redis 服务器发送命令并返回收到的答复
（3）redis.Bool(reply interface{}, err error)：将命令返回转为布尔值
（4）redis.Bytes(reply interface{}, err error)：将命令返回转为 Bytes
（5）redis.Strings(reply interface{}, err error)：将命令返回转为 []string
*/

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisPool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// HSet a key/value
func HSet(key string, data interface{}, time int) error {
	conn := RedisPool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("HSet", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
