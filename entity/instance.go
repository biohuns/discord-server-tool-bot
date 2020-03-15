package entity

// InstanceService インスタンスサービス
type InstanceService interface {
	Start() error
	Stop() error
	Status() (*InstanceInfo, error)
}

// InstanceInfo インスタンス情報
type InstanceInfo struct {
	Name   string         `json:"name"`
	Status InstanceStatus `json:"status"`
}

// InstanceStatus インスタンス状態
type InstanceStatus int

const (
	InstanceStatusUnknown InstanceStatus = iota
	InstanceStatusPending
	InstanceStatusRunning
	InstanceStatusStopping
	InstanceStatusStopped
)

func (is InstanceStatus) String() string {
	switch is {
	case InstanceStatusPending:
		return "起動処理中"
	case InstanceStatusRunning:
		return "起動"
	case InstanceStatusStopping:
		return "停止処理中"
	case InstanceStatusStopped:
		return "停止"
	default:
		return "不明"
	}
}
