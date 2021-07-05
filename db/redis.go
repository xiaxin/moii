package db

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xiaxin/moii/db/redis/lock"
	"github.com/xiaxin/moii/log"
)

var (
	ErrRedisClosed = errors.New("redis is already closed")
	ErrKeyNotFound = errors.New("key '%s' doesn't found")
)

type Redis struct {
	pool   *redis.Pool
	config *RedisConfig
}

type RedisConfig struct {
	Host      string `yaml:"host"`
	Password  string `yaml:"password"`
	DB        int    `yaml:"db"`
	MaxIdle   int    `yaml:"max_idle"`
	MaxActive int    `yaml:"max_active"`
	Wait      bool   `yaml:"wait"`
	Prefix    string `yaml:"prefix"`
}

func NewRedis(config *RedisConfig) *Redis {

	if nil == config {
		return nil
	}

	pool := NewRedisPool(config)

	r := &Redis{
		pool:   pool,
		config: config,
	}

	return r
}

func NewRedisPool(config *RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   config.MaxIdle,
		MaxActive: config.MaxActive,

		IdleTimeout: 5 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Host)
			if err != nil {
				return nil, err
			}

			if len(config.Password) > 0 {
				if _, err = c.Do("AUTH", config.Password); err != nil {
					log.Info(err)
					c.Close()
					return nil, err
				}
			}

			if _, err = c.Do("SELECT", config.DB); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		// 无法获得可用连接时，将会暂时阻塞
		Wait: config.Wait,
	}
}

func (r *Redis) Pool() *redis.Pool {
	return r.pool
}

func (r *Redis) PingPong() (bool, error) {
	c := r.pool.Get()
	defer c.Close()
	msg, err := c.Do("PING")

	if err != nil || msg == nil {
		return false, err
	}

	return msg == "PONG", nil
}

func (r *Redis) CloseConnection() error {
	if r.pool != nil {
		return r.pool.Close()
	}
	return ErrRedisClosed
}

func (r *Redis) Send(command string, args ...interface{}) (err error) {
	return r.pool.Get().Send(command, args...)
}

func (r *Redis) Do(command string, args ...interface{}) (reply interface{}, err error) {
	c := r.pool.Get()

	defer c.Close()

	if err := c.Err(); nil != c.Err() {
		return nil, err
	}

	return c.Do(command, args...)
}

func (r *Redis) Lock(resource, token string, timeout int, fc func(ok bool) error) error {
	lock, ok, err := lock.TryLock(r.pool.Get(), resource, token, timeout)

	if err != nil {
		return err
	}

	ren := fc(ok)

	if ok {
		defer lock.Unlock()
	}

	return ren
}
