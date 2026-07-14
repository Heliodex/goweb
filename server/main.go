package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Heliodex/goweb/shared"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("../index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(file)
	})

	// files
	http.HandleFunc("/main.wasm", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("../client/main.wasm")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/wasm")
		w.Write(file)
	})

	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("../client/wasm_exec.js")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/javascript")
		w.Write(file)
	})

	HandleRemoteFunc(shared.ThingFunc)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
