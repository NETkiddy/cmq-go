package cmq_go

import (
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
	"bytes"
	"errors"
	"net"
)

const (
	DEFAULT_HTTP_TIMEOUT = 3000 //ms
)

type CMQHttp struct {
	isKeepAlive bool
	conn        *http.Client
}

var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          500,
	MaxIdleConnsPerHost:   100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func NewCMQHttp() *CMQHttp {
	return &CMQHttp{
		isKeepAlive: true,
		conn: &http.Client{
			Timeout:   DEFAULT_HTTP_TIMEOUT * time.Millisecond,
			Transport: DefaultTransport,
		},
	}

}

func (this *CMQHttp) setProxy(proxyUrlStr string) (err error) {
	if proxyUrlStr == "" {
		return
	}
	proxyUrl, err := url.Parse(proxyUrlStr)
	if err != nil {
		return
	}
	transport, err := this.getTransport()
	if err != nil {
		return
	}
	transport.Proxy = http.ProxyURL(proxyUrl)
	return
}

func (this *CMQHttp) unsetProxy() (err error) {
	transport, err := this.getTransport()
	if err != nil {
		return
	}
	transport.Proxy = nil
	return
}

func (this *CMQHttp) getTransport() (*http.Transport, error) {
	if this.conn.Transport == nil {
		this.SetTransport(DefaultTransport)
	}

	if transport, ok := this.conn.Transport.(*http.Transport); ok {
		return transport, nil
	}
	return nil, errors.New("transport is not an *http.Transport instance")
}

func (this *CMQHttp) SetTransport(transport http.RoundTripper) {
	this.conn.Transport = transport
}

func (this *CMQHttp) request(method, urlStr, reqStr string, userTimeout int) (result string, err error) {
	timeout := DEFAULT_HTTP_TIMEOUT
	if userTimeout >= 0 {
		timeout += userTimeout
	}
	this.conn.Timeout = time.Duration(timeout) * time.Millisecond

	req, err := http.NewRequest(method, urlStr, bytes.NewReader([]byte(reqStr)))
	if err != nil {
		return "", fmt.Errorf("make http req error %v", err)
	}
	resp, err := this.conn.Do(req)
	if err != nil {
		return "", fmt.Errorf("http error  %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http error code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read http resp body error %v", err)
	}
	result = string(body)
	return
}
