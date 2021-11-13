package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/req_counter_service/config"
)

type Database struct {
	Client *redis.Client
}

var (
	Ctx = context.TODO()
)

func newClient() (*Database, error) {
	configurations, err_config := config.ConfigService()
	if err_config != nil {
		fmt.Println(fmt.Errorf("Failed to connect to client with error: %s", err_config.Error()))
		return nil, err_config
	}

	db_address := configurations.Container_db_host + ":" + configurations.Container_db_port

	var redisClient = redis.NewClient(&redis.Options{
		Addr:     db_address,
		Password: "",
		DB:       0,
	})

	if err := redisClient.Ping(Ctx).Err(); err != nil {
		return nil, err
	}

	return &Database{
		Client: redisClient,
	}, nil
}

/* @Parameter: Integer (a number that is the total number of requests in the cluster).
   @Returns : void
*/
func SetData(payload int) {
	var ctx = context.TODO()

	db, err_client := newClient()
	if err_client != nil {
		fmt.Println(fmt.Errorf("Failed to connect to client with error: %s", err_client.Error()))
	}

	db.Client.Set(ctx, "total_count", payload, 0)
}

/* @Parameter: no parameters.
   @Returns : String (the total number of requests in the cluster until now.)
*/
func GetData() string {
	var ctx = context.TODO()

	db, err_client := newClient()
	if err_client != nil {
		fmt.Println(fmt.Errorf("Failed to connect to client with error: %s", err_client.Error()))
	}

	total_count, err := db.Client.Get(ctx, "total_count").Result()
	if err != nil {
		fmt.Println("Error with the get total_count")
		panic(err)
	}
	return total_count
}
