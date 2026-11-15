package shared

type (
	Data interface {
		Validate() error
	}

	RemoteFunc[ReqType, ResType Data] struct {
		Name string
		Callback func(ReqType) ResType
	}
)
