//go:build js && wasm

package main

import "github.com/Heliodex/goweb/shared"

func main() {
	println("Hello from the client!")

	res, err := Invoke(shared.ThingFunc, shared.Thing{A: "Hello", B: 42})
	if err != nil {
		panic(err)
	}

	println("Response from server:", res.A, res.B)

	// doc := js.Global().Get("document")
	// h1 := doc.Call("createElement", "h1")
	// h1.Set("textContent", "Hello from Go WASM!")
	// h1.Set("style", "color: white")
	// doc.Get("body").Call("appendChild", h1)

	dom := Dom{
		Body: []Element{
			el("h1", Attrs{
				"style": "color: white",
			}, []Element{
				text("Hello from Go WASM!"),
			}),
		},
	}

	dom.Render()
}
