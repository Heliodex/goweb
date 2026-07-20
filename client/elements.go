//go:build js && wasm

package main

import "syscall/js"

type Attrs map[string]any

type Element interface {
	Render() js.Value
}

type Elements []Element

type TagElement struct {
	name     string
	attrs    Attrs
	children Elements
}

func (te TagElement) Children(children ...Element) TagElement {
	te.children = children
	return te
}

func (te TagElement) Attr(key string, value any) TagElement {
	if te.attrs == nil {
		te.attrs = make(Attrs, 1)
	}
	te.attrs[key] = value
	return te
}

func (te TagElement) Render() js.Value {
	doc := js.Global().Get("document")
	el := doc.Call("createElement", te.name)
	for k, v := range te.attrs {
		el.Set(k, v)
	}
	for _, child := range te.children {
		el.Call("appendChild", renderElement(child))
	}
	return el
}

type TextNode struct {
	Text string
}

func (tn TextNode) Render() js.Value {
	doc := js.Global().Get("document")
	return doc.Call("createTextNode", tn.Text)
}

func e(name string) TagElement {
	return TagElement{name: name}
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

func (cte *ComputedElement) Render() js.Value {
	node := TagElement(cte.compute(nilNotifier)).Render()
	node.Call("setAttribute", "data-dynamic-id", cte.ID)

	doc := js.Global().Get("document")
	// if an element with the same ID already exists, replace it
	if existing := doc.Call("querySelector", "[data-dynamic-id='"+cte.ID+"']"); existing.Truthy() {
		existing.Call("replaceWith", node)
	}

	return node
}
