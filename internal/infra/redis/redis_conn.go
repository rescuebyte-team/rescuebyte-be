package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

func NewRedisConnWithDB(db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     env.AppEnv.RedisAddress,
		Password: env.AppEnv.RedisPassword,
		DB:       db,
	})

	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Warn(log.Fields{
			"error": err.Error(),
		}, fmt.Sprintf("[Redis] Failed to connect to Redis on %s (db %d)", env.AppEnv.RedisAddress, db))
	} else {
		log.Info(log.Fields{
			"ping": ping,
		}, fmt.Sprintf("[Redis] Successfully connected to Redis (db %d)", db))
	}

	return client
}
