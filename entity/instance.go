package entity

// InstanceService インスタンスサービス
type InstanceService interface {
	Start() error
	Stop() error
	GetStatus() (*InstanceStatus, error)
	GetAndCacheStatus() (*InstanceStatus, error)
	GetCachedStatus() (*InstanceStatus, error)
}

// InstanceStatus インスタンスステータス
type InstanceStatus struct {
	Name            string             `json:"name"`
	StatusCode      InstanceStatusCode `json:"status"`
	IsStatusChanged bool               `json:"is_status_changed"`
}

// InstanceStatusCode インスタンスステータスコード
type InstanceStatusCode int

const (
	InstanceStatusUnknown InstanceStatusCode = iota
	InstanceStatusPending
	InstanceStatusRunning
	InstanceStatusStopping
	InstanceStatusStopped
)

func (is InstanceStatusCode) String() string {
	switch is {
	case InstanceStatusPending:
		return "Start Processing"
	case InstanceStatusRunning:
		return "Running"
	case InstanceStatusStopping:
		return "Stop Processing"
	case InstanceStatusStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}
