package shared

type Thing struct {
	A string
	B int
}

var ThingFunc = RemoteFunc[Thing, Thing]{
	Name: "ThingFunc",
}

func (t Thing) Validate() error {
	return nil // always valid stub
}
