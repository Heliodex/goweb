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

type StaticElement interface {
	Element
	StaticElement()
}

type TagElement struct {
	Name     string
	Attrs    Attrs
	Children []Element
}

func (TagElement) Element()       {}
func (TagElement) StaticElement() {}

type TextNode struct {
	Text string
}

func (TextNode) Element()       {}
func (TextNode) StaticElement() {}

func el(name string, attrs Attrs, children []Element) TagElement {
	return TagElement{name, attrs, children}
}

type Dynamic[T any] interface {
	ID() string
	Function() T
}

type DynamicElement struct {
	id       string
	function func() TagElement
}

func (de DynamicElement) ID() string {
	return de.id
}

func (de DynamicElement) Function() TagElement {
	return de.function()
}

func (DynamicElement) Element() {}

// todo: dynamicelement constructor

func text(text string) StaticElement {
	return TextNode{text}
}

type Dom struct {
	Head, Body []Element
}

func renderTagElement(se TagElement) js.Value {
	doc := js.Global().Get("document")
	el := doc.Call("createElement", se.Name)
	for k, v := range se.Attrs {
		el.Set(k, v)
	}
	for _, child := range se.Children {
		el.Call("appendChild", renderElement(child))
	}
	return el
}

func renderTextNode(tn TextNode) js.Value {
	doc := js.Global().Get("document")
	return doc.Call("createTextNode", tn.Text)
}

func renderStaticElement(se StaticElement) js.Value {
	switch v := se.(type) {
	case TagElement:
		return renderTagElement(v)
	case TextNode:
		return renderTextNode(v)
	}

	panic("unknown static element type")
}

func renderDynamicElement(de DynamicElement) js.Value {
	node := renderStaticElement(de.function())
	node.Call("setAttribute", "data-dynamic-id", de.id)

	doc := js.Global().Get("document")
	// if an element with the same ID already exists, replace it
	if existing := doc.Call("querySelector", "[data-dynamic-id='"+de.id+"']"); existing.Truthy() {
		existing.Call("replaceWith", node)
	}

	return node
}

func renderElement(e Element) js.Value {
	switch v := e.(type) {
	case StaticElement:
		return renderStaticElement(v)
	case DynamicElement:
		return renderDynamicElement(v)
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

type Value[T any] struct {
	ID           string
	Value        T
	dependencies map[*Point[T]]struct{}
}

func Val[T any](v T) *Value[T] {
	return &Value[T]{
		Value:        v,
		dependencies: make(map[*Point[T]]struct{}),
	}
}

func (v *Value[T]) Set(newValue T) {
	v.Value = newValue

	// trigger update
	for p := range v.dependencies {
		renderDynamicElement(*p.de)
	}
}

func Peek[T any](v *Value[T]) T {
	return v.Value
}

type Point[T any] struct {
	de *Dynamic[T]
}

func Use[T any](p *Point[T], v *Value[T]) T {
	v.dependencies[p] = struct{}{}
	return v.Value
}

func Dyn[T any](f func(p *Point[T]) T) Dynamic[T] {
	p := &Point[T]{}

	de := Dynamic[T]{
		ID: randomStringID(),
		// Value: f(p),
		Function: func() T {
			return f(p)
		},
	}
	p.de = &de

	return de
}
