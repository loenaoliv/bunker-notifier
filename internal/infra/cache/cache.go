package cache

import (
	"encoding/json"

	"github.com/coocood/freecache"
)

type FreeCache interface {
	Get(key string, r interface{}) error
	Set(key string, obj interface{}, expireSeconds int) error
	Delete(key string) bool
}

type freeCacheImpl struct {
	cache *freecache.Cache
}

func NewCache(cacheSize int) FreeCache {
	return &freeCacheImpl{
		cache: freecache.NewCache(cacheSize),
	}
}

func (fc *freeCacheImpl) Get(key string, r interface{}) error {
	object, err := fc.cache.Get([]byte(key))
	if err != nil {
		return err
	}

	if err = json.Unmarshal(object, &r); err != nil {
		return err
	}

	return nil
}

func (fc *freeCacheImpl) Set(key string, obj interface{}, expireSeconds int) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	if err := fc.cache.Set([]byte(key), b, expireSeconds); err != nil {
		return err
	}

	return nil
}

func (fc *freeCacheImpl) Delete(key string) bool {
	return fc.cache.Del([]byte(key))
}
