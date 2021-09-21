package cache

import (
	"github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

var Cache *cache.Cache

func init() {
	once := sync.Once{}
	once.Do(func() {
		log.Println("Creating cache..")
		if Cache == nil {
			Cache = cache.New(5*time.Minute, 10*time.Minute)
		}
	})
}

func Memory() *cache.Cache {
	return Cache
}
