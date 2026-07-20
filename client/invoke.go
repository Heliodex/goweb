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

func Invoke[ReqType, ResType shared.Data](f shared.RemoteFunc[ReqType, ResType], arg ReqType, callback func(ResType)) (err error) {
	if err = arg.Validate(); err != nil {
		return fmt.Errorf("invalid argument: %w", err)
	}

	var w bytes.Buffer
	if err = gob.NewEncoder(&w).Encode(arg); err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	hres, err := http.Post(url+"/api/"+f.Name, "application/octet-stream", &w)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer hres.Body.Close()

	if hres.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned bad status: %s", hres.Status)
	}

	var res ResType
	if err = gob.NewDecoder(hres.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if err = res.Validate(); err != nil {
		return fmt.Errorf("invalid response: %w", err)
	}

	fmt.Println("response! sending...")
	callback(res)
	return
}
