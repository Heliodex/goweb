//go:build js && wasm

package main

import (
	"math/rand"
	"syscall/js"
)

const idChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomStringID() string {
	id := make([]byte, 8)
	for i := range 8 {
		id[i] = idChars[rand.Intn(len(idChars))]
	}

	return string(id)
}

func text(text string) Element {
	return TextNode{text}
}

type Dom struct {
	Head, Body Elements
}

type CustomNotifier struct {
	notifyFunc func()
}

func (cn *CustomNotifier) Notify() {
	cn.notifyFunc()
}

func (cn *CustomNotifier) AddDependent(Notifiable)  {}
func (cn *CustomNotifier) AddDependency(Notifiable) {}

func renderElement(se Element) js.Value {
	switch v := se.(type) {
	case TagElement, TextNode:
		return v.Render()
	case *ComputedElement:
		notify := &CustomNotifier{}
		notify.notifyFunc = func() {
			println("Notifying dependents of computed element with ID:", v.ID)
			v.Render()
		}

		v.AddDependent(notify)

		return v.Render()
	}

	panic("unknown element type")
}

func (d Dom) Render() {
	doc := js.Global().Get("document")
	for _, elem := range d.Head {
		node := renderElement(elem)
		doc.Get("head").Call("appendChild", node)
	}
	for _, elem := range d.Body {
		node := renderElement(elem)
		doc.Get("body").Call("appendChild", node)
	}

	select {} // keep program open for event handlers
}

func MakeFunc(fn func()) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		fn()
		return nil
	})
}

type Notifiable interface {
	Notify()
	AddDependent(Notifiable)
	AddDependency(Notifiable)
}

type Deps struct {
	dependencies, dependents map[Notifiable]struct{}
}

func MakeDeps() Deps {
	return Deps{
		dependencies: make(map[Notifiable]struct{}),
		dependents:   make(map[Notifiable]struct{}),
	}
}

func (d *Deps) Dependencies() map[Notifiable]struct{} {
	return d.dependencies
}

func (d *Deps) Dependents() map[Notifiable]struct{} {
	return d.dependents
}

type Notifier func(Notifiable)

func nilNotifier(Notifiable) {}
