package cache

import "github.com/MoScenix/mes/common/filestore/base"

func NewStore(actual base.Store) (CacheStore, error) {
	return NewLocalCacheStore(actual)
}
