package ctl

import (
	"encoding/json"
	"github.com/decred/dcrd/dcrjson/v3"
	"time"
)

const rpcVersion = "1.0"

func CallGet(result interface{}, url, method string, params ...interface{}) error {
	id := time.Now().Unix()

	req, err := dcrjson.NewRequest(rpcVersion, id, method, params)
	if err != nil {
		return err
	}
	marshalledJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}
	rsp, err := rpcCallByUrl(url, marshalledJSON)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}
	return nil
}

func CallPost(result interface{}, url, method string, params ...interface{}) (string, error) {
	id := time.Now().Unix()

	req, err := dcrjson.NewRequest(rpcVersion, id, method, params)
	if err != nil {
		return "", err
	}
	marshalledJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	rsp, err := rpcCallByUrl(url, marshalledJSON)
	if err != nil {
		return "", err
	}
	return string(rsp), nil
}
