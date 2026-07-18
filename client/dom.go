//go:build js && wasm

package main

import (
	"syscall/js"
)

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

func el(name string, attrs Attrs, children []Element) Element {
	return TagElement{name, attrs, children}
}

type DynamicElement struct {
	Value StaticElement
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

func renderElement(e Element) js.Value {
	switch v := e.(type) {
	case StaticElement:
		return renderStaticElement(v)
	case DynamicElement:
		return renderStaticElement(v.Value)
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

type GenericValue[T any] interface {
	Set(T)
}

type Value[T any] struct {
	Value T
}

func Val[T any](v T) *Value[T] {
	return &Value[T]{Value: v}
}

func (v *Value[T]) Set(newValue T) {
	v.Value = newValue
}

func Peek[T any](v *Value[T]) T {
	return v.Value
}

type Point struct {
	dependencies map[any]struct{} // *Value
}

func Use[T any](p *Point, v *Value[T]) T {
	p.dependencies[v] = struct{}{}
	return v.Value
}

func Dynamic(f func(p *Point) StaticElement) DynamicElement {
	p := &Point{
		dependencies: make(map[any]struct{}),
	}

	return DynamicElement{
		Value: f(p),
	}
}
