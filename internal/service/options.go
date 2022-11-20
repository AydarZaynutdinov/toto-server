package service

import (
	"toto-server/internal/repository"

	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
)

type Options struct {
	Log           *zap.Logger
	SkuConfigRepo repository.ISkuConfig
	RedisClient   *redis.Client
}

type OptionFunc func(*Options)

func WithLogger(log *zap.Logger) OptionFunc {
	return func(s *Options) {
		s.Log = log
	}
}

func WithSkuConfigRepo(repo repository.ISkuConfig) OptionFunc {
	return func(s *Options) {
		s.SkuConfigRepo = repo
	}
}

func WithRedisClient(redisClient *redis.Client) OptionFunc {
	return func(s *Options) {
		s.RedisClient = redisClient
	}
}
