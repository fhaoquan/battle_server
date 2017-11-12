package result_cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var(
	default_result_cache *cache.Cache;
)

func init(){
	default_result_cache=cache.New(time.Hour,time.Hour);
}

func CacheResult(key string,value interface{}){
	default_result_cache.Set(key,value,cache.DefaultExpiration);
}

func GetResult(key string)(interface{}, bool) {
	return default_result_cache.Get(key);
}