//go:build js && wasm

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"syscall/js"

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

func Fetch(url, method string, body *bytes.Buffer) (*http.Response, error) {
	f := js.Global().Get("fetch")
	if !f.Truthy() {
		return nil, fmt.Errorf("fetch is not available in this environment")
	}

	headers := js.Global().Get("Headers").New()
	headers.Call("append", "Content-Type", "application/octet-stream")

	options := js.ValueOf(map[string]any{
		"method":  method,
		"headers": headers,
		"body":    body.Bytes(),
	})

	promise := f.Invoke(url, options)
	then := promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		return args[0]
	}))

	catch := then.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
		err := args[0]
		fmt.Println("Fetch error:", err)
		return nil
	}))

	response := catch.Call("await")
	if !response.Truthy() {
		return nil, fmt.Errorf("fetch failed")
	}

	status := response.Get("status").Int()
	if status != http.StatusOK {
		return nil, fmt.Errorf("server returned bad status: %d", status)
	}

	arrayBuffer := response.Call("arrayBuffer").Call("await")
	data := js.Global().Get("Uint8Array").New(arrayBuffer)
	buf := make([]byte, data.Get("length").Int())
	js.CopyBytesToGo(buf, data)

	return &http.Response{
		StatusCode: status,
		Body:       http.NoBody,
	}, nil
}

func InvokeFetch[ReqType, ResType shared.Data](f shared.RemoteFunc[ReqType, ResType], arg ReqType) (res ResType, err error) {
	if err = arg.Validate(); err != nil {
		return res, fmt.Errorf("invalid argument: %w", err)
	}

	var w bytes.Buffer
	if err = gob.NewEncoder(&w).Encode(arg); err != nil {
		return res, fmt.Errorf("failed to encode request: %w", err)
	}

	hres, err := Fetch(url+"/api/"+f.Name, "POST", &w)
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
