//go:build js && wasm

package main

import "fmt"

// Values based on Computeds
type Value[T any] struct {
	Computed[T]
}

func NewValue[T any](initial T) *Value[T] {
	v := &Value[T]{}

	v.Computed = Computed[T]{
		ID:   randomStringID(),
		Deps: MakeDeps(),
		compute: func(n Notifier) T {
			n(v)
			return initial
		},
	}

	return v
}

func (v *Value[T]) Set(newValue T) {
	v.Computed.compute = func(Notifier) T {
		return newValue
	}

	fmt.Println("Value set called, notifying dependents")
	for dep := range v.dependents {
		fmt.Println("notifying", dep)
		dep.Notify()
	}
}
