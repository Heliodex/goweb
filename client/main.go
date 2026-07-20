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
			e("h1").
				Attr("style", "color: white").
				Children(
					text("Hello from Go WASM!"),
				),

			NewComputedElement(func(n Notifier) TagElement {
				return e("p").
					Attr("style", "color: white").
					Children(
						text("You have clicked the button " + strconv.Itoa(num.Use(n)) + " times."),
					)
			}),

			NewComputedElement(func(n Notifier) TagElement {
				return e("p").
					Attr("style", "color: white").
					Children(
						text("Quadruple that equals " + strconv.Itoa(quadruple.Use(n)) + " times."),
					)
			}),

			e("button").
				On("click", func() {
					num.Set(num.Peek() + 1)
				}).
				Children(
					text("Click me"),
				),
		},
	}

	dom.Render()
}
