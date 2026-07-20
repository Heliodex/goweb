package shared

type Data interface {
	Validate() error
}

type RemoteFunc[ReqType, ResType Data] struct {
	Name     string
	Callback func(ReqType) ResType
}
