package jsonrpc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/bytedance/sonic"
	"net/http"
	"os"
	"time"
)

type RPCResult[T any] struct {
	Result T      `json:"result"`
	Id     string `json:"id"`
	Error  *any   `json:"error"`
}

type RPCRequest struct {
	Version string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

func CallRPC[T any](version, method string, params any) (*RPCResult[T], error) {
	client := http.Client{Timeout: 5 * time.Second}

	request := &RPCRequest{
		Version: version,
		Id:      fmt.Sprintf("litepay_%s_%d", method, time.Now().Unix()),
		Method:  method,
		Params:  params,
	}

	buffer := new(bytes.Buffer)
	if err := sonic.ConfigDefault.NewEncoder(buffer).Encode(request); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", os.Getenv("LTC_RPC_HOST"), buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("LTC_RPC_USER"), os.Getenv("LTC_RPC_PASSWORD"))))))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result RPCResult[T]
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		errBuffer := new(bytes.Buffer)
		if err := sonic.ConfigDefault.NewEncoder(errBuffer).Encode(result.Error); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("error: %s", errBuffer.String())
	}

	return &result, nil
}
