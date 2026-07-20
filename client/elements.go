//go:build js && wasm

package main

import "syscall/js"

type Attrs map[string]any

type Element interface {
	Render() js.Value
}

type Elements []Element

type TagElement struct {
	Name     string
	Attrs    Attrs
	Children Elements
}

func (te TagElement) Render() js.Value {
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

type TextNode struct {
	Text string
}

func (tn TextNode) Render() js.Value {
	doc := js.Global().Get("document")
	return doc.Call("createTextNode", tn.Text)
}

func el(name string, attrs Attrs, children Elements) TagElement {
	return TagElement{name, attrs, children}
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
