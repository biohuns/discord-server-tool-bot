package entity

type (
	// CacheService キャッシュサービス
	CacheService interface {
		Get(key string, v interface{}) error
		Set(key string, value interface{}) error
	}
)
