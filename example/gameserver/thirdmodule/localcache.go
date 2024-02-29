package thirdmodule

import (
	"context"

	"github.com/allegro/bigcache/v3"
)

var (
	LocalCache *bigcache.BigCache
)

func InitLocalCache() {
	var err error
	LocalCache, err = bigcache.New(context.Background(), bigcache.DefaultConfig(-1))
	if err != nil {
		panic(err)
	}
}
