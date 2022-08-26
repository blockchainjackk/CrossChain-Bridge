// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ctl

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/decred/dcrd/dcrjson/v3"
)

var (
	dcrnConfigs   []*params.DcrnConfig
	dcrnConfigMap map[string]*params.DcrnConfig
)

func initArgs() {
	if dcrnConfigMap == nil {
		cfg := params.GetConfig()
		dcrnConfigs = cfg.DcrnConfigs
		dcrnConfigMap = make(map[string]*params.DcrnConfig)
		for _, dcrnConfig := range dcrnConfigs {
			dcrnConfigMap[dcrnConfig.IP] = dcrnConfig
		}
	}

}

// notls方式
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

// tls
func newHTTPClient2(config *ConnConfig) (*http.Client, error) {
	var dial func(network, addr string) (net.Conn, error)
	var tlsConfig *tls.Config

	// Configure TLS if needed.
	if !config.DisableTLS {
		if len(config.Certificates) > 0 {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(config.Certificates)
			tlsConfig = &tls.Config{
				RootCAs: pool,
			}
		}
	}
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

func rpcCallByUrl(urlStr string, marshalledJSON []byte) ([]byte, error) {
	initArgs()
	bodyReader := bytes.NewReader(marshalledJSON)
	httpRequest, err := http.NewRequest("POST", urlStr, bodyReader)
	if err != nil {
		return nil, err
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")

	// Create the new HTTP client that is configured according to the user-
	// specified options and submit the request.
	connConfig, err := getConnConfig(urlStr)
	if err != nil {
		return nil, err
	}
	// Configure basic access authorization.
	httpRequest.SetBasicAuth(connConfig.User, connConfig.Pass)

	// httpClient, err := newHTTPClient()
	httpClient, err := newHTTPClient2(connConfig)
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

// ConnConfig describes the connection configuration parameters for the client.
// This
type ConnConfig struct {
	// Host is the IP address and port of the RPC server you want to connect
	// to.
	Host string

	// Endpoint is the websocket endpoint on the RPC server.  This is
	// typically "ws".
	Endpoint string

	// User is the username to use to authenticate to the RPC server.
	User string

	// Pass is the passphrase to use to authenticate to the RPC server.
	Pass string

	// DisableTLS specifies whether transport layer security should be
	// disabled.  It is recommended to always use TLS if the RPC server
	// supports it as otherwise your username and password is sent across
	// the wire in cleartext.
	DisableTLS bool

	// Certificates are the bytes for a PEM-encoded certificate chain used
	// for the TLS connection.  It has no effect if the DisableTLS parameter
	// is true.
	Certificates []byte
}

func getConnConfig(urlStr string) (*ConnConfig, error) {
	ip, port, err := parseIPAndPort(urlStr)
	if err != nil {
		return nil, err
	}
	path, err := getCertPathByPort(ip, port)
	if err != nil {
		return nil, err
	}
	serverCert, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := dcrnConfigMap[ip]
	return &ConnConfig{
		User:         config.Rpcuser,
		Pass:         config.Rpcpassword,
		DisableTLS:   false,
		Certificates: serverCert,
	}, nil
}

// 通过url获取IP与端口号
func parseIPAndPort(urlStr string) (string, int, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", 0, err
	}
	host := u.Host
	ho := strings.Split(host, ":")
	ip := ho[0]
	port, err := strconv.Atoi(ho[1])
	if err != nil {
		return "", 0, err
	}
	return ip, port, nil
}

// 通过端口号确认证书路径
func getCertPathByPort(ip string, port int) (string, error) {
	dcrnConfig := dcrnConfigMap[ip]
	if dcrnConfig == nil {
		return "", errors.New("ip:" + ip + "no config")
	}
	var path string
	switch port {
	case 9109: //mainnet
		fallthrough
	case 19109: //testnet
		fallthrough
	case 19556: //simnet
		fallthrough
	case dcrnConfig.DcrndPort: //配置自定义
		path = dcrnConfig.DcrndRpcCertPath //dcrnd证书路径
	case 9110: //mainnet
		fallthrough
	case 19110: //testnet
		fallthrough
	case 19557: //simnet
		fallthrough
	case dcrnConfig.WalletPort: //配置自定义
		path = dcrnConfig.WalletRpcCertPath //wallet证书路径
	default:
		return "", errors.New("port error")
	}
	return path, nil
}
