package entity

// CacheService キャッシュサービス
type CacheService interface {
	Get(key CacheKey) (interface{}, error)
	Set(key CacheKey, value interface{}) error
	Dump() CacheDumpList
}

// CacheKey キャッシュキー
type CacheKey int

const (
	ServerStatusKey CacheKey = iota
	InstanceStatusKey
)

func (c CacheKey) String() string {
	switch c {
	case ServerStatusKey:
		return "ServerStatus"
	case InstanceStatusKey:
		return "InstanceStatus"
	default:
		return "Unknown"
	}
}

// CacheDumpList キャッシュダンプデータリスト
type CacheDumpList []*CacheDump

// CacheDumpElement キャッシュダンプデータ
type CacheDump struct {
	Key   CacheKey    `json:"key"`
	Value interface{} `json:"value"`
}
