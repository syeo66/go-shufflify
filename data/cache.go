package data

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var cacheStore = cache.New(5*time.Minute, 10*time.Minute)
