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

func (d Dom) Render() {
	doc := js.Global().Get("document")
}
