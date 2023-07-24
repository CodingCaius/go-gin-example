//redis工具包

package gredis

import (
	"encoding/json"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/gomodule/redigo/redis"
)

// 全局变量，整个程序中只有一个连接池实例
var RedisConn *redis.Pool

// 初始化 Redis 连接池
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		//当连接池需要新的连接时，会调用这个函数来创建并返回一个新的Redis连接
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		//用于检查从连接池中获取的连接是否仍然有效
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// 将数据存储到Redis中，并设置过期时间
func Set(key string, data interface{}, time int) error {
	//从之前设置的Redis连接池中获取一个连接
	/*
	调用 Get() 方法时，连接池会检查是否有可用的空闲连接。如果有，它会返回一个现有的连接；如果没有，它会调用之前设置的 Dial 函数来创建一个新的 Redis 连接，然后将这个新连接返回
	*/
	conn := RedisConn.Get()
	//确保在函数结束后关闭Redis连接，以便将连接放回连接池中，避免资源泄漏
	defer conn.Close()

	//将data数据进行JSON序列化，将其转换为字节切片value
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	/*
	使用Redis连接执行SET命令，将键名key和序列化后的数据value存储到Redis中。
	使用conn.Do()方法可以向Redis发送各种命令，并获取执行结果
	*/
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

//检查指定的键名key是否存在于Redis中
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

//从Redis中根据指定的键名key获取数据
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}


//用于模糊匹配指定的键名key，并删除匹配到的所有键对应的数据
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	//redis.Strings()函数用于将命令返回的结果转换为字符串切片，其中包含所有匹配的键名
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