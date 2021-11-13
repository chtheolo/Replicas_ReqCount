package redis

import (
	"context"
	"fmt"
	"strconv"

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
func IncrGet_total_count() string {
	var ctx = context.TODO()

	db, err_client := newClient()
	if err_client != nil {
		fmt.Println(fmt.Errorf("Failed to connect to client with error: %s", err_client.Error()))
	}
	pipe := db.Client.TxPipeline()

	total_count := pipe.Incr(ctx, "total_count")

	_, err_pipe := pipe.Exec(ctx)
	if err_pipe != nil {
		panic(err_pipe)
	}

	return strconv.FormatInt(total_count.Val(), 10)
}
