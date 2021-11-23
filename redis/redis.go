package redis

import (
	"github.com/gomodule/redigo/redis"
	"otpapp-native/config"
	"time"
)

var (
	conn      redis.Conn
	redisPool redis.Pool
)

func Init(cfg config.RedisCfg) (redis.Conn, error) {
	redisPool = redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", cfg.URL+":"+cfg.Port)
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("AUTH", cfg.Password)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 60 * time.Second,
	}

	conn = redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}

	return conn, err
}

func String(command string, params ...interface{}) (string, error) {
	conn := redisPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do(command, params...))
	if err == redis.ErrNil {
		return res, nil
	}
	return res, err
}

func Strings(command string, params ...interface{}) ([]string, error) {
	conn := redisPool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do(command, params...))
	return res, err
}

func Int(command string, params ...interface{}) (int, error) {
	conn := redisPool.Get()
	defer conn.Close()

	res, err := redis.Int(conn.Do(command, params...))
	return res, err
}
