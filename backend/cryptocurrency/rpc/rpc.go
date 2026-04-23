package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type Response[T any] struct {
	Result T `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func Call[T any](url, method string, params []any) (T, error) {
	var zero T
	body, err := json.Marshal(Request{JSONRPC: "1.0", Method: method, Params: params, ID: 1})
	if err != nil {
		return zero, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	var result Response[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return zero, err
	}
	if result.Error != nil {
		return zero, fmt.Errorf("rpc %s: %s (code %d)", method, result.Error.Message, result.Error.Code)
	}
	return result.Result, nil
}
