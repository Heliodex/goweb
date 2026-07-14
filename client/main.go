package main

import "github.com/Heliodex/goweb/shared"

func main() {
	println("Hello from the client!")

	res, err := Invoke(shared.ThingFunc, shared.Thing{A: "Hello", B: 42})
	if err != nil {
		panic(err)
	}

	println("Response from server:", res.A, res.B)
}
