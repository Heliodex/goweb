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

type Attrs map[string]any

type Element interface {
	Element()
}

type TagElement struct {
	Name     string
	Attrs    Attrs
	Children []Element
}

func (TagElement) Element() {}

type TextNode struct {
	Text string
}

func (TextNode) Element() {}

func el(name string, attrs Attrs, children []Element) TagElement {
	return TagElement{name, attrs, children}
}

// todo: dynamicelement constructor

func text(text string) Element {
	return TextNode{text}
}

type Dom struct {
	Head, Body []Element
}

func renderTagElement(te TagElement) js.Value {
	doc := js.Global().Get("document")
	el := doc.Call("createElement", te.Name)
	for k, v := range te.Attrs {
		el.Set(k, v)
	}
	for _, child := range te.Children {
		el.Call("appendChild", renderElement(child))
	}
	return el
}

func renderTextNode(tn TextNode) js.Value {
	doc := js.Global().Get("document")
	return doc.Call("createTextNode", tn.Text)
}

func renderTagElementComputed(cte *Computed[TagElement]) js.Value {
	node := renderTagElement(cte.compute(nilNotifier))
	node.Call("setAttribute", "data-dynamic-id", cte.ID)

	doc := js.Global().Get("document")
	// if an element with the same ID already exists, replace it
	if existing := doc.Call("querySelector", "[data-dynamic-id='"+cte.ID+"']"); existing.Truthy() {
		existing.Call("replaceWith", node)
	}

	return node
}

func renderElement(se Element) js.Value {
	switch v := se.(type) {
	case TagElement:
		return renderTagElement(v)
	case TextNode:
		return renderTextNode(v)
	case *Computed[TagElement]:
		return renderTagElementComputed(v)
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

type State interface {
	Dependencies() map[State]struct{}
	Dependents() map[State]struct{}
	Compute() any
}

type Deps struct {
	dependencies, dependents map[State]struct{}
}

func MakeDeps() Deps {
	return Deps{
		dependencies: make(map[State]struct{}),
		dependents:   make(map[State]struct{}),
	}
}

func (d *Deps) Dependencies() map[State]struct{} {
	return d.dependencies
}

func (d *Deps) Dependents() map[State]struct{} {
	return d.dependents
}

type notifier func(State)

func nilNotifier(State) {}

type Compute[T any] func(notifier) T

type Computed[T any] struct {
	ID string
	Deps
	compute Compute[T]
}

func NewComputed[T any](compute Compute[T]) *Computed[T] {
	deps := MakeDeps()

	notifier := func(s State) {
		deps.dependencies[s] = struct{}{}
	}
	compute(notifier)

	return &Computed[T]{
		ID:      randomStringID(),
		Deps:    deps,
		compute: compute,
	}
}

func (c *Computed[T]) Peek() T {
	return c.compute(nilNotifier)
}

func (c *Computed[T]) Use(n notifier) T {
	return c.compute(n)
}

func (*Computed[T]) Element() {} // MAYBE?

// Values based on Computeds
type Value[T any] struct {
	Computed[T]
}

func NewValue[T any](initial T) *Value[T] {
	return &Value[T]{
		Computed: *NewComputed(func(notifier) T {
			return initial // Initial T
		}),
	}
}

func (v *Value[T]) Set(newValue T) {
	v.Computed.compute = func(notifier) T {
		return newValue
	}

	// TODO: notify dependents
}
