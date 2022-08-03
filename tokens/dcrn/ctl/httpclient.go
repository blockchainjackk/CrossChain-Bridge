// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ctl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/decred/dcrd/dcrjson/v3"
	"io/ioutil"
	"net"
	"net/http"
)

const urlDefault = "http://127.0.0.1:19557"

var (
	url         string
	rpcuser     string
	rpcpassword string
)

func initArgs() {
	rpcuser = "root"
	rpcpassword = "root"
}

func newHTTPClient() (*http.Client, error) {
	var dial func(network, addr string) (net.Conn, error)
	var tlsConfig *tls.Config
	// Create and return the new HTTP client potentially configured with a
	// proxy and TLS.
	client := http.Client{
		Transport: &http.Transport{
			Dial:            dial,
			TLSClientConfig: tlsConfig,
		},
	}
	return &client, nil
}

func rpcCallByUrl(url string, marshalledJSON []byte) ([]byte, error) {

	bodyReader := bytes.NewReader(marshalledJSON)
	httpRequest, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")

	// Configure basic access authorization.
	// httpRequest.SetBasicAuth(rpcuser, rpcpassword)

	// Create the new HTTP client that is configured according to the user-
	// specified options and submit the request.
	httpClient, err := newHTTPClient()
	if err != nil {
		return nil, err
	}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	// Read the raw bytes and close the response.
	respBytes, err := ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		err = fmt.Errorf("error reading json reply: %w", err)
		return nil, err
	}

	// Handle unsuccessful HTTP responses
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		// Generate a standard error to return if the server body is
		// empty.  This should not happen very often, but it's better
		// than showing nothing in case the target server has a poor
		// implementation.
		if len(respBytes) == 0 {
			return nil, fmt.Errorf("%d %s", httpResponse.StatusCode,
				http.StatusText(httpResponse.StatusCode))
		}
		return nil, fmt.Errorf("%s", respBytes)
	}

	fmt.Println(string(respBytes))

	// Unmarshal the response.
	var resp dcrjson.Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func rpcCall(marshalledJSON []byte) ([]byte, error) {
	initArgs()
	return rpcCallByUrl(urlDefault, marshalledJSON)
}
