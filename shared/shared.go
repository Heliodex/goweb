package shared

type Data interface {
	Validate() error
}

type RemoteFunc[ReqType Data, ResType Data] struct {
	Name     string
	Callback func(ReqType) ResType
}
