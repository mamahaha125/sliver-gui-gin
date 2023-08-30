package cache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	ca   *cache.Cache
	once sync.Once
)

func Instance() *cache.Cache {
	once.Do(func() {
		ca = cache.New(1*time.Minute, 3*time.Minute)
	})
	return ca
}
