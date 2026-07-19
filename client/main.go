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

	num := NewValue(0)
	double := NewComputed(func(n Notifier) int {
		return num.Use(n) * 2
	})
	quadruple := NewComputed(func(n Notifier) int {
		return double.Use(n) * 2
	})

	dom := Dom{
		Body: Elements{
			el("h1", Attrs{
				"style": "color: white",
			}, Elements{
				text("Hello from Go WASM!"),
			}),

			NewComputedElement(func(n Notifier) TagElement {
				println("Dynamic function called! num is", num.Peek())
				unum := num.Use(n)

				return el("p", Attrs{
					"style": "color: white",
				}, Elements{
					text("You have clicked the button " + strconv.Itoa(unum) + " times."),
				})
			}),

			NewComputedElement(func(n Notifier) TagElement {
				println("Dynamic function called! num is", num.Peek())
				uquadruple := quadruple.Use(n)

				return el("p", Attrs{
					"style": "color: white",
				}, Elements{
					text("Quadruple that equals " + strconv.Itoa(uquadruple) + " times."),
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
