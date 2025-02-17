// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"context"

	"github.com/xzl8028/xenia-server/einterfaces"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/utils"
)

const (
	REACTION_CACHE_SIZE = 20000
	REACTION_CACHE_SEC  = 30 * 60

	ROLE_CACHE_SIZE = 20000
	ROLE_CACHE_SEC  = 30 * 60

	SCHEME_CACHE_SIZE = 20000
	SCHEME_CACHE_SEC  = 30 * 60

	GROUP_CACHE_SIZE = 20000
	GROUP_CACHE_SEC  = 30 * 60

	CLEAR_CACHE_MESSAGE_DATA = ""
)

type LocalCacheSupplier struct {
	next          LayeredStoreSupplier
	reactionCache *utils.Cache
	roleCache     *utils.Cache
	schemeCache   *utils.Cache
	metrics       einterfaces.MetricsInterface
	cluster       einterfaces.ClusterInterface
}

// Caching Interface
type ObjectCache interface {
	AddWithExpiresInSecs(key, value interface{}, expireAtSecs int64)
	AddWithDefaultExpires(key, value interface{})
	Purge()
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
	Len() int
	Name() string
	GetInvalidateClusterEvent() string
}

func NewLocalCacheSupplier(metrics einterfaces.MetricsInterface, cluster einterfaces.ClusterInterface) *LocalCacheSupplier {
	supplier := &LocalCacheSupplier{
		reactionCache: utils.NewLruWithParams(REACTION_CACHE_SIZE, "Reaction", REACTION_CACHE_SEC, model.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_REACTIONS),
		roleCache:     utils.NewLruWithParams(ROLE_CACHE_SIZE, "Role", ROLE_CACHE_SEC, model.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_ROLES),
		schemeCache:   utils.NewLruWithParams(SCHEME_CACHE_SIZE, "Scheme", SCHEME_CACHE_SEC, model.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_SCHEMES),
		metrics:       metrics,
		cluster:       cluster,
	}

	if cluster != nil {
		cluster.RegisterClusterMessageHandler(model.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_REACTIONS, supplier.handleClusterInvalidateReaction)
		cluster.RegisterClusterMessageHandler(model.CLUSTER_EVENT_INVALIDATE_CACHE_FOR_ROLES, supplier.handleClusterInvalidateRole)
	}

	return supplier
}

func (s *LocalCacheSupplier) SetChainNext(next LayeredStoreSupplier) {
	s.next = next
}

func (s *LocalCacheSupplier) Next() LayeredStoreSupplier {
	return s.next
}

func (s *LocalCacheSupplier) doStandardReadCache(ctx context.Context, cache ObjectCache, key string, hints ...LayeredStoreHint) *LayeredStoreSupplierResult {
	if hintsContains(hints, LSH_NO_CACHE) {
		if s.metrics != nil {
			s.metrics.IncrementMemCacheMissCounter(cache.Name())
		}
		return nil
	}

	if cacheItem, ok := cache.Get(key); ok {
		if s.metrics != nil {
			s.metrics.IncrementMemCacheHitCounter(cache.Name())
		}
		result := NewSupplierResult()
		result.Data = cacheItem
		return result
	}

	if s.metrics != nil {
		s.metrics.IncrementMemCacheMissCounter(cache.Name())
	}

	return nil
}

func (s *LocalCacheSupplier) doStandardAddToCache(ctx context.Context, cache ObjectCache, key string, result *LayeredStoreSupplierResult, hints ...LayeredStoreHint) {
	if result.Err == nil && result.Data != nil {
		cache.AddWithDefaultExpires(key, result.Data)
	}
}

func (s *LocalCacheSupplier) doInvalidateCacheCluster(cache ObjectCache, key string) {
	cache.Remove(key)
	if s.cluster != nil {
		msg := &model.ClusterMessage{
			Event:    cache.GetInvalidateClusterEvent(),
			SendType: model.CLUSTER_SEND_BEST_EFFORT,
			Data:     key,
		}
		s.cluster.SendClusterMessage(msg)
	}
}

func (s *LocalCacheSupplier) doClearCacheCluster(cache ObjectCache) {
	cache.Purge()
	if s.cluster != nil {
		msg := &model.ClusterMessage{
			Event:    cache.GetInvalidateClusterEvent(),
			SendType: model.CLUSTER_SEND_BEST_EFFORT,
			Data:     CLEAR_CACHE_MESSAGE_DATA,
		}
		s.cluster.SendClusterMessage(msg)
	}
}

func (s *LocalCacheSupplier) Invalidate() {
	s.doClearCacheCluster(s.reactionCache)
	s.doClearCacheCluster(s.roleCache)
	s.doClearCacheCluster(s.schemeCache)
}
