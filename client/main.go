//go:build js && wasm

package main

import (
	"strconv"

	"github.com/Heliodex/goweb/shared"
)

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

	num := Val(0)

	dom := Dom{
		Body: []Element{
			el("h1", Attrs{
				"style": "color: white",
			}, []Element{
				text("Hello from Go WASM!"),
			}),

			Dynamic(func(p *Point) TagElement {
				println("Dynamic function called! num is", Peek(num))
				n := Use(p, num)

				return el("p", Attrs{
					"style": "color: white",
				}, []Element{
					text("You have clicked the button " + strconv.Itoa(n) + " times."),
				})
			}),

			el("button", Attrs{
				"onclick": MakeFunc(func() {
					num.Set(Peek(num) + 1)
					println("Button clicked, num is now", Peek(num))
				}),
			}, []Element{
				text("Click me"),
			}),
		},
	}

	dom.Render()
}
