package db

import (
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestRedis(t *testing.T) {
	r := NewRedis(&RedisConfig{
		Host:     "192.168.1.206:6379",
		Password: "880728",
		DB:       0,
	})

	res1, err1 := redis.String(r.Do("set", "a", "b"))

	if nil != err1 || res1 != "OK" {
		t.Error("set failed")
	}

	res2, err2 := redis.String(r.Do("get", "a"))
	if nil != err2 || res2 != "b" {
		t.Error("get failed")
	}
}
