package data

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var CacheStore = cache.New(5*time.Minute, 10*time.Minute)
