package utils

import (
	"log"
	"sync"

	lru "github.com/hashicorp/golang-lru"
)

var (
    StaticCache *lru.Cache
    cacheOnce   sync.Once
)

var (
    DynamicCache *lru.Cache
    dynamicCacheOnce sync.Once
)


func InitStaticCache() {
    cacheOnce.Do(func() {
        var err error
        StaticCache, err = lru.New(128) // Cache up to 128 items
        if err != nil {
            log.Fatalf("Error initializing static cache: %v", err)
        }
    })
}


func InitDynamicCache() {
    dynamicCacheOnce.Do(func() {
        var err error
        DynamicCache, err = lru.New(128) // Adjust size as needed
        if err != nil {
            log.Fatalf("Error initializing dynamic cache: %v", err)
        }
    })
}
