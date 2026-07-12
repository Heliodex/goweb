package shared

import (
	"encoding/gob"
	"io"
)

type Thing struct {
	A string
	B int
}

func (t Thing) Serialise(w io.Writer) error {
	return gob.NewEncoder(w).Encode(t)
}

func DeserialiseThing(r io.Reader) (t Thing, err error) {
	return t, gob.NewDecoder(r).Decode(&t)
}
