package gredis

import (
	"encoding/json"
	"fmt"
	"gbins/gconf"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdb = &redis.Client{}
	dsn = ""
	serviceName = "redis"
)

func init()  {
	selectRedis()
	err := setRedis()
	if err != nil {
		panic(fmt.Sprintf("redis is wrong : ",err.Error))
	}
}

func NewRedis() *redis.Client {
	_ = setRedis()
	return rdb
}

func selectRedis()  {
	if gconf.GConfByClientName(serviceName).Target != ""{
		serviceName = gconf.GConfByClientName(serviceName).Target
	}
	return
}

func setRedis() error {
	if dsn != gconf.GetRVal(serviceName){
		dsn = gconf.GetRVal(serviceName)
		redisOpt := defaultRedisOpt()
		err := json.Unmarshal([]byte(dsn) , redisOpt)
		if err != nil {
			return err
		}
		r := redis.NewClient(redisOpt)
		initHook(r)
		rdb = r
	}
	return nil
}

func initHook(r *redis.Client)  {
	r.AddHook(&gredisLogHook{})
	return
}

func defaultRedisOpt() *redis.Options {
	return &redis.Options{
		DialTimeout : 800 * time.Millisecond,
		PoolSize : 1000,
	}
}
