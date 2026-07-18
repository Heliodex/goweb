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

type DynamicElement struct {
	ID       string
	Function func() TagElement
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
	node := renderStaticElement(de.Function())
	node.Call("setAttribute", "data-dynamic-id", de.ID)

	doc := js.Global().Get("document")
	// if an element with the same ID already exists, replace it
	if existing := doc.Call("querySelector", "[data-dynamic-id='"+de.ID+"']"); existing.Truthy() {
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
	Value        T
	dependencies map[*Point]struct{}
}

func Val[T any](v T) *Value[T] {
	return &Value[T]{
		Value:        v,
		dependencies: make(map[*Point]struct{}),
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

type Point struct {
	de *DynamicElement
}

func Use[T any](p *Point, v *Value[T]) T {
	v.dependencies[p] = struct{}{}
	return v.Value
}

func Dynamic(f func(p *Point) TagElement) DynamicElement {
	p := &Point{}

	de := DynamicElement{
		ID: randomStringID(),
		// Value: f(p),
		Function: func() TagElement {
			return f(p)
		},
	}
	p.de = &de

	return de
}
