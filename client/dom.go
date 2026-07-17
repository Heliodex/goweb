//go:build js && wasm

package main

import "syscall/js"

type Attrs map[string]string

type Element interface {
	Element()
}

type StaticElement struct {
	Name     string
	Attrs    Attrs
	Children []Element
}

func (StaticElement) Element() {}

func el(name string, attrs Attrs, children []Element) Element {
	return StaticElement{name, attrs, children}
}

type DynamicElement struct {
	Value  StaticElement
	Update chan struct{}
}

func (DynamicElement) Element() {}

// todo: dynamicelement constructor

type TextNode struct {
	Text string
}

func (TextNode) Element() {}

func text(text string) Element {
	return TextNode{text}
}

type Dom struct {
	Head, Body []Element
}

func renderStaticElement(se StaticElement) js.Value {
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

func renderElement(e Element) js.Value {
	switch v := e.(type) {
	case StaticElement:
		return renderStaticElement(v)
	case TextNode:
		return renderTextNode(v)
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
}
