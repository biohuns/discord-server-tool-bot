package cache

import (
	"encoding/json"
	"sync"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

// Service キャッシュサービス
type Service struct {
	s *sync.Map
}

var (
	serviceInstance *Service
	once            sync.Once
)

// ProvideService サービス返却
func ProvideService() (entity.CacheService, error) {
	var err error

	once.Do(func() {
		serviceInstance = &Service{s: new(sync.Map)}
	})

	if serviceInstance == nil {
		err = xerrors.New("service is not provided")
	}

	if err != nil {
		return nil, xerrors.Errorf("provide service error: %w", err)
	}

	return serviceInstance, nil
}

// Get キャッシュから取り出す
func (s Service) Get(key string, v interface{}) error {
	value, ok := s.s.Load(key)
	if !ok {
		return xerrors.New("key not found")
	}

	b, ok := value.([]byte)
	if !ok {
		return xerrors.New("failed to assertion")
	}

	if err := json.Unmarshal(b, v); err != nil {
		return xerrors.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

// Set キャッシュに保存
func (s Service) Set(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return xerrors.Errorf("failed to marshal json: %w", err)
	}

	s.s.Store(key, b)
	return nil
}
