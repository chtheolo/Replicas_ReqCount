package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/req_counter_service/config"
)

// Database is struct variable for a redis Client
type Database struct {
	Client *redis.Client
}

// Generates a new Redis database client
func newClient(ctx context.Context) (*Database, error) {
	
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

// IncrGetTotalCount is a function that is responsible for making the INCR command transaction in Redis
// and returns back the total request counter in string form. A Redis error is returned in case of
//  something went wrong with the transaction/pipeline or a cancel error that is generated from the root ctx.
func IncrGetTotalCount(ctx context.Context) (string, error) {
	
	select {
	case <- ctx.Done():
		err := ctx.Err()
		return "", err
	default:
		db, errClient := newClient(ctx)
		if errClient != nil {
			return "", errClient
		}
		pipe := db.Client.TxPipeline()

		totalCount := pipe.Incr(ctx, "total_count")

		_, errPipe := pipe.Exec(ctx)
		if errPipe != nil {
			return "", errPipe
		}

		return strconv.FormatInt(totalCount.Val(), 10), nil
	}
}