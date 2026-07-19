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

	num := NewValue(0)
	double := NewComputed(func(n notifier) int {
		udouble := num.Use(n)

		return udouble * 2
	})

	dom := Dom{
		Body: Elements{
			el("h1", Attrs{
				"style": "color: white",
			}, Elements{
				text("Hello from Go WASM!"),
			}),

			NewComputedElement(func(n notifier) TagElement {
				println("Dynamic function called! num is", num.Peek())
				unum := num.Use(n)

				return el("p", Attrs{
					"style": "color: white",
				}, Elements{
					text("You have clicked the button " + strconv.Itoa(unum) + " times."),
				})
			}),

			NewComputedElement(func(n notifier) TagElement {
				println("Dynamic function called! num is", num.Peek())
				udouble := double.Use(n)

				return el("p", Attrs{
					"style": "color: white",
				}, Elements{
					text("Double that equals " + strconv.Itoa(udouble) + " times."),
				})
			}),

			el("button", Attrs{
				"onclick": MakeFunc(func() {
					num.Set(num.Peek() + 1)
					println("Button clicked, num is now", num.Peek())
				}),
			}, Elements{
				text("Click me"),
			}),
		},
	}

	dom.Render()
}
