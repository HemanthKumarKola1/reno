package main

import (
	"net/http"

	"github.com/go-redis/redis"
	"router.com/repo"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.ListenAndServe(":8080", route(repo.NewRedisClient(rdb)))
}
