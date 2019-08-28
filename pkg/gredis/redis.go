package gredis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/gin-blog/pkg/setting"
)

var RedisConn *redis.Pool

func Setup() error {

	//创建一个Redis连接池（设置了最大活跃连接数，最大空闲连接数，空闲时间）
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,

		//设置redis监听地址，redis密码
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
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func Set(key string, data interface{}, time int) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	reply, err := redis.Bool(conn.Do("SET", key, value))
	conn.Do("EXPIRE", key, time)

	return reply, err
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

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

//模糊删除
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	//根据参数，获取redis中存储的包含参数的key
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	//遍历删除模糊查询得到的key
	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}