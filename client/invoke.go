//go:build js && wasm

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/Heliodex/goweb/shared"
)

const url = "http://localhost:8080"

func Invoke[ReqType, ResType shared.Data](f shared.RemoteFunc[ReqType, ResType], arg ReqType) (res ResType, err error) {
	if err = arg.Validate(); err != nil {
		return res, fmt.Errorf("invalid argument: %w", err)
	}

	var w bytes.Buffer
	if err = gob.NewEncoder(&w).Encode(arg); err != nil {
		return res, fmt.Errorf("failed to encode request: %w", err)
	}

	hres, err := http.Post(url+"/api/"+f.Name, "application/octet-stream", &w)
	if err != nil {
		return res, fmt.Errorf("failed to send request: %w", err)
	}
	defer hres.Body.Close()

	if hres.StatusCode != http.StatusOK {
		return res, fmt.Errorf("server returned bad status: %s", hres.Status)
	}

	var r bytes.Buffer
	if _, err = r.ReadFrom(hres.Body); err != nil {
		return res, fmt.Errorf("failed to read response body: %w", err)
	}

	if err = gob.NewDecoder(&r).Decode(&res); err != nil {
		return res, fmt.Errorf("failed to decode response: %w", err)
	}

	if err = res.Validate(); err != nil {
		return res, fmt.Errorf("invalid response: %w", err)
	}

	return res, nil
}
