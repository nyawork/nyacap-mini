package inits

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"nya-captcha/config"
	g "nya-captcha/global"
	"time"
)

func Redis() error {
	// Parse connect string
	redisConfig, err := redis.ParseURL(config.Config.System.Redis)
	if err != nil {
		return fmt.Errorf("格式化 redis 连接 url: %w", err)
	}

	// Connect to server
	g.Redis = redis.NewClient(redisConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Try connection
	err = g.Redis.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("连接到 redis: %w", err)
	}

	return nil
}
