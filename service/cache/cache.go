package cache

import (
	"sync"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

// Service キャッシュサービス
type Service struct {
	s *sync.Map
}

// Get キャッシュから取り出す
func (s *Service) Get(key entity.CacheKey) (interface{}, error) {
	value, ok := s.s.Load(key)
	if !ok {
		return nil, xerrors.New("specified key cannot be found")
	}

	return value, nil
}

// Set キャッシュに保存
func (s *Service) Set(key entity.CacheKey, value interface{}) error {
	s.s.Store(key, value)
	return nil
}

// Dump キャッシュをダンプする
func (s *Service) Dump() entity.CacheDumpList {
	list := make(entity.CacheDumpList, 0)
	s.s.Range(func(key, value interface{}) bool {
		k := key.(entity.CacheKey)
		list = append(list, &entity.CacheDump{
			Key:   k,
			Value: value,
		})
		return true
	})

	return list
}

var (
	shared *Service
	once   sync.Once
)

// ProvideService サービス返却
func ProvideService() (entity.CacheService, error) {
	once.Do(func() {
		shared = &Service{s: new(sync.Map)}
	})

	if shared == nil {
		return nil, xerrors.New("service is not provided")
	}

	return shared, nil
}
