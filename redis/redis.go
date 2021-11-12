package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "db:6379",
	Password: "",
	DB:       0,
})

func SetData(payload int) {
	// var ctx = context.Background()
	var ctx = context.TODO()

	redisClient.Set(ctx, "total_count", payload, 0)
}

func GetData() string {
	// var ctx = context.Background()
	var ctx = context.TODO()

	total_count, err := redisClient.Get(ctx, "total_count").Result()
	if err != nil {
		fmt.Println("Error with the get total_count")
		panic(err)
	}
	return total_count
}
