package repo

import (
	"Graduation-Project/config"
	redislib "github.com/go-redis/redis/v8"
	"time"
)

type RedisDB struct {
	*redislib.Client
}
type RedisClientOpt struct {
	Addr         string
	Password     string
	DB           int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
	PoolSize     int
	MinIdleConns int
}

func InitRedisDefault() (*RedisDB, error) {
	opt := RedisClientOpt{}
	opt.Addr = config.Conf.RedisAddress
	opt.Password = config.Conf.RedisPassword
	opt.DB = 0
	//opt.ReadTimeout = time.Duration(10)
	//opt.WriteTimeout = time.Duration(10)
	//opt.DialTimeout = time.Duration(10)
	//opt.PoolSize = 100
	//opt.MinIdleConns = 10
	rdb, err := InitRedis(opt)
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
func InitRedis(opt RedisClientOpt) (*RedisDB, error) {
	rdb := redislib.NewClient(&redislib.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	redis := &RedisDB{rdb}

	return redis, nil
}
