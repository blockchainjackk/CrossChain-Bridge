package ctl

import (
	"encoding/json"
	"time"

	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/decred/dcrd/dcrjson/v3"
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
		log.Warn("rpc err", "url", url, "method", method, "err", err)
		return err
	}
	//log.Info("rpc success", "url", url, "method", method, "result", string(rsp))
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}
	return nil
}

func CallPost(url, method string, params ...interface{}) (string, error) {
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
