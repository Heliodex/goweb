//go:build js && wasm

package main

import (
	"fmt"
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

type Elements []Element

type TagElement struct {
	Name     string
	Attrs    Attrs
	Children Elements
}

func (TagElement) Element() {}

type TextNode struct {
	Text string
}

func (TextNode) Element() {}

func el(name string, attrs Attrs, children Elements) TagElement {
	return TagElement{name, attrs, children}
}

// todo: dynamicelement constructor

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
	case *ComputedElement:
		notify := &CustomNotifier{}
		notify.notifyFunc = func() {
			println("Notifying dependents of computed element with ID:", v.ID)
			renderTagElementComputed(&v.Computed)
		}

		v.AddDependent(notify)

		return renderTagElementComputed(&v.Computed)
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

type notifier func(Notifiable)

func nilNotifier(Notifiable) {}

type Compute[T any] func(notifier) T

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

func (c *Computed[T]) Use(n notifier) T {
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

type ComputedElement struct {
	Computed[TagElement]
}

func NewComputedElement(compute Compute[TagElement]) *ComputedElement {
	c := NewComputed(compute)

	return &ComputedElement{
		Computed: *c,
	}
}

func (ComputedElement) Element() {}

// Values based on Computeds
type Value[T any] struct {
	Computed[T]
}

func NewValue[T any](initial T) *Value[T] {
	v := &Value[T]{}

	v.Computed = Computed[T]{
		ID:   randomStringID(),
		Deps: MakeDeps(),
		compute: func(n notifier) T {
			n(v)
			return initial
		},
	}

	return v
}

func (v *Value[T]) Set(newValue T) {
	v.Computed.compute = func(notifier) T {
		return newValue
	}

	fmt.Println("Value set called, notifying dependents")
	for dep := range v.dependents {
		fmt.Println("notifying", dep)
		dep.Notify()
	}
}
