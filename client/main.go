package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	println("Hello from the client!")

	res, err := http.Post("/api", "application/octet-stream", nil)
	if err != nil {
		panic(err)
	}

	println("Response status:", res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response body: ", body)

	defer res.Body.Close()
}
