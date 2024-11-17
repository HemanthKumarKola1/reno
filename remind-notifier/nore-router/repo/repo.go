package repo

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	rc *redis.Client
}

func NewRedisClient(rc *redis.Client) *RedisClient {
	return &RedisClient{rc}
}

func (r *RedisClient) Validate(uName string, pwd string) error {
	val, err := r.rc.Get(uName).Result()
	if err != nil {
		return err
	}
	if !verifyPassword(pwd, val) {
		return fmt.Errorf("wrong password")
	}
	return nil
}

func verifyPassword(password string, storedHash string) bool {      
	return password == storedHash
}
