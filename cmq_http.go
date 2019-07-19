package cmq_go

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
	"encoding/json"
)

type CMQHttp struct {
	timeout     int
	isKeepAlive bool
	conn        *http.Client
	debugMode   bool
}

func NewCMQHttp() *CMQHttp {
	return &CMQHttp{
		timeout:     10000,
		isKeepAlive: true,
		conn:        nil,
		debugMode:   false,
	}
}

//Debug open the CMQHttp debug mode
func (c *CMQHttp) SetDebug(b bool) {
	c.debugMode = b
}

func (c *CMQHttp) DebugMode() bool {
	return c.debugMode
}

type simpleHTTPReq struct {
	Method           string
	Proto            string // "HTTP/1.0"
	TransferEncoding []string
	Host             string
	RemoteAddr       string
	RequestURI       string
}

func genSimpleHTTPReq(r *http.Request) *simpleHTTPReq {
	return &simpleHTTPReq{
		Method:           r.Method,
		Proto:            r.Proto,
		TransferEncoding: r.TransferEncoding,
		Host:             r.Host,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
	}
}

func printSimpleReqInfo(r *http.Request) {
	date := time.Now().Format("2019-07-19 15:11::00")
	info, _ := json.Marshal(genSimpleHTTPReq(r))
	fmt.Printf("time: %s | info: %s", date, string(info))
}

func (this *CMQHttp) request(method, urlStr, reqStr, proxyUrlStr string, userTimeout int) (result string, err error) {
	var client *http.Client
	timeout := 0
	if userTimeout >= 0 {
		timeout = userTimeout
	}
	keepalive := 0
	if this.isKeepAlive {
		keepalive = 30
	}

	if proxyUrlStr == "" {
		unproxyTransport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: time.Duration(keepalive) * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		client = &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: unproxyTransport,
		}
	} else {
		proxyUrl, err := url.Parse(proxyUrlStr)
		if err != nil {
			panic(err)
		}
		proxyTransport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: time.Duration(keepalive) * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		client = &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: proxyTransport,
		}
	}

	req, err := http.NewRequest(method, urlStr, bytes.NewReader([]byte(reqStr)))
	if err != nil {
		return "", fmt.Errorf("make http req error %v", err)
	}
	if this.DebugMode() {
		printSimpleReqInfo(req)
	}
	resp, err := client.Do(req)
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
