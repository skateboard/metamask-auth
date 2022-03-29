package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Connection = Cache{
	Client: cache.New(5 * time.Minute, 10 * time.Minute),
	DefaultExpiration: 5 * time.Minute,
}

type Cache struct {
	Client *cache.Cache
	DefaultExpiration time.Duration
}