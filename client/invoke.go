package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/Heliodex/goweb/shared"
)

const url = "http://localhost:8080"

func Invoke[ReqType shared.Data, ResType shared.Data](f shared.RemoteFunc[ReqType, ResType], arg ReqType) (res ResType, err error) {
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

	if err = gob.NewDecoder(hres.Body).Decode(&res); err != nil {
		return res, fmt.Errorf("failed to decode response: %w", err)
	}

	return
}
