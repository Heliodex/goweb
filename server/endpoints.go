package main

import (
	"encoding/gob"
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

func TransformRemoteFunc[ReqType shared.Data, ResType shared.Data](f shared.RemoteFunc[ReqType, ResType]) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var in ReqType
		if gob.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
			return
		}

		out := f.Callback(in)

		if gob.NewEncoder(w).Encode(out) != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
