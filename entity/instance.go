package entity

type (
	// InstanceService インスタンスサービス
	InstanceService interface {
		Start() error
		Stop() error
		Status() (InstanceInfo, error)
	}

	// InstanceInfo インスタンス情報
	InstanceInfo interface {
		GetStatus() string
	}
)
