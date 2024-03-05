package utils

import "github.com/go-redis/redis/v8"

type Configuration struct {
	EnvName string
}
var Param Configuration
var RedisConnection *redis.Client