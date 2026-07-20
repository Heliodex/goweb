//go:build js && wasm

package main

import (
	"fmt"
	"strconv"

	"github.com/Heliodex/goweb/shared"
)

func main() {
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

			NewComputedElement(func(n Notifier) TagElement {
				nquadruple := quadruple.Use(n)

				fmt.Println("nquadruple:", nquadruple, "calling Invoke with ThingFunc")

				responseChan := make(chan shared.Thing, 1)
				if err := Invoke(shared.ThingFunc, shared.Thing{A: "Hello", B: nquadruple}, func(res shared.Thing) {
					fmt.Println("sending response...")
					responseChan <- res
					fmt.Println("sent!")
				}); err != nil {
					panic(err)
				}

				fmt.Println("Request finished")
				res := <-responseChan
				fmt.Println("Received response from server:", res)

				return e("p").
					Attr("style", "color: white").
					Children(
						text("Updated from the server that's " + strconv.Itoa(res.B) + " times."),
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
