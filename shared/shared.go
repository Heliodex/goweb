package shared

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
)

type Data interface {
	Validate() error
}

type RemoteFunc[ReqType Data, ResType Data] struct {
	Name     string
	Callback func(ReqType) ResType
}

const url = "http://localhost:8080"

func (f RemoteFunc[ReqType, ResType]) invokeServer(arg ReqType) (res ResType, err error) {
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

func (f RemoteFunc[ReqType, ResType]) Invoke(arg ReqType) (ResType, error) {
	if f.Callback != nil {
		return f.Callback(arg), nil
	}

	return f.invokeServer(arg)
}
