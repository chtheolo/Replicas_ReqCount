package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/req_counter_service/config"
)

/*Database is struct variable for a redis Client*/
type Database struct {
	Client *redis.Client
}

var (
	ctx = context.Background()
)

func newClient() (*Database, error) {
	configurations, errInit := config.Initializer()
	if errInit != nil {
		fmt.Println(fmt.Errorf("Failed to connect to client with error: %s", errInit.Error()))
		return nil, errInit
	}

	dbAddress := fmt.Sprintf("%s:%s",configurations.ContainerDBhost, configurations.ContainerDBport)

	var redisClient = redis.NewClient(&redis.Options{
		Addr:     dbAddress,
		Password: "",
		DB:       0,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Database{
		Client: redisClient,
	}, nil
}

/*IncrGetTotalCount ...
@Parameter: Integer (a number that is the total number of requests in the cluster).
@Returns : string OR error
*/
func IncrGetTotalCount() (string, error) {

	db, errClient := newClient()
	if errClient != nil {
		fmt.Println(fmt.Errorf("Failed to connect to client with error: %s", errClient.Error()))
	}
	pipe := db.Client.TxPipeline()

	totalCount := pipe.Incr(ctx, "total_count")

	_, errPipe := pipe.Exec(ctx)
	if errPipe != nil {
		return "", errPipe
	}

	return strconv.FormatInt(totalCount.Val(), 10), nil
}
