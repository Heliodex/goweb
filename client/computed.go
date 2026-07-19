//go:build js && wasm

package main

import "fmt"

type Compute[T any] func(Notifier) T

type Computed[T any] struct {
	ID string
	Deps
	compute Compute[T]
}

func NewComputed[T any](compute Compute[T]) *Computed[T] {
	deps := MakeDeps()

	c := &Computed[T]{
		ID:      randomStringID(),
		Deps:    deps,
		compute: compute,
	}

	notifier := func(s Notifiable) {
		c.AddDependency(s)
		s.AddDependent(c)
		fmt.Println("notifier called, added dependency:", s)
	}
	compute(notifier)

	return c
}

func (c *Computed[T]) Peek() T {
	return c.compute(nilNotifier)
}

func (c *Computed[T]) Use(n Notifier) T {
	return c.compute(n)
}

func (c *Computed[T]) Notify() {
	fmt.Println("Computed Notify called, notifying dependents")
	for dep := range c.dependents {
		dep.Notify()
	}
}

func (c *Computed[T]) AddDependent(dep Notifiable) {
	c.dependents[dep] = struct{}{}
	fmt.Println("Added dependent:", dep)
}

func (c *Computed[T]) AddDependency(dep Notifiable) {
	c.dependencies[dep] = struct{}{}
	fmt.Println("Added dependency:", dep)
}
