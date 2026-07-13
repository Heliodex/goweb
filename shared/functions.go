package shared

type Thing struct {
	A string
	B int
}

func (t Thing) Validate() error {
	return nil // always valid stub
}

var ThingFunc = RemoteFunc[Thing, Thing]{
	Name: "ThingFunc",
}
