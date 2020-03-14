package entity

const (
	StatusUnknown InstanceStatus = iota
	StatusProvisioning
	StatusRepairing
	StatusRunning
	StatusStaging
	StatusStopped
	StatusStopping
	StatusSuspended
	StatusSuspending
	StatusTerminated
)

// InstanceService インスタンスサービス
type InstanceService interface {
	Start() error
	Stop() error
	Status() (InstanceInfo, error)
}

// InstanceInfo インスタンス情報
type InstanceInfo interface {
	GetStatus() InstanceStatus
	GetFormattedStatus() string
}

// InstanceStatus インスタンス状態
type InstanceStatus int

func (is InstanceStatus) String() string {
	switch is {
	case StatusProvisioning:
		return "リソース割当中"
	case StatusRepairing:
		return "修復中"
	case StatusRunning:
		return "起動"
	case StatusStaging:
		return "起動準備中"
	case StatusStopped:
		return "停止"
	case StatusStopping:
		return "停止準備中"
	case StatusSuspended:
		return "休止"
	case StatusSuspending:
		return "休止準備中"
	case StatusTerminated:
		return "終了"
	default:
		return "不明"
	}
}
