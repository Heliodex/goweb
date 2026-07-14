package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"

	"github.com/Heliodex/goweb/shared"
)

func init() {
	shared.ThingFunc.Callback = func(t shared.Thing) shared.Thing {
		return shared.Thing{
			A: t.A + " processed",
			B: t.B + 1,
		}
	}
}

func HandleRemoteFunc[ReqType shared.Data, ResType shared.Data](f shared.RemoteFunc[ReqType, ResType]) {
	http.HandleFunc("POST /api/"+f.Name, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var in ReqType
		if gob.NewDecoder(bytes.NewBuffer(body)).Decode(&in) != nil {
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
			return
		}

		out := f.Callback(in)

		if gob.NewEncoder(w).Encode(out) != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
	})

	fmt.Println("Registered remote function:", f.Name)
}
