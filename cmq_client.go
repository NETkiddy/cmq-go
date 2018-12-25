package cmq_go

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
	"net/url"
	"strconv"
	"sort"
)

const (
	CURRENT_VERSION = "SDK_GO_1.3"
)

type CMQClient struct {
	Endpoint   string
	Path       string
	SecretId   string
	SecretKey  string
	Method     string
	SignMethod string
	Proxy      string
	CmqHttp    *CMQHttp
}

func NewCMQClient(endpoint, path, secretId, secretKey, method string) *CMQClient {
	return &CMQClient{
		Endpoint:   endpoint,
		Path:       path,
		SecretId:   secretId,
		SecretKey:  secretKey,
		Method:     method,
		SignMethod: "sha1",
		CmqHttp:    NewCMQHttp(),
	}
}

func (this *CMQClient) setSignMethod(signMethod string) (err error) {
	if signMethod != "sha1" && signMethod != "sha256" {
		err = fmt.Errorf("Only support sha1 or sha256 now")
		return
	} else {
		this.SignMethod = signMethod
	}
	return
}

func (this *CMQClient) setProxy(proxyUrl string) {
	this.Proxy = proxyUrl
	return
}

func (this *CMQClient) unsetProxy() {
	this.Proxy = ""
	return
}

func (this *CMQClient) call(action string, param map[string]string) (resp string, err error) {
	param["Action"] = action
	param["Nonce"] = strconv.Itoa(rand.Int())
	param["SecretId"] = this.SecretId
	param["Timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	param["RequestClient"] = CURRENT_VERSION
	if this.SignMethod == "sha256" {
		param["SignatureMethod"] = "HmacSHA256"
	} else {
		param["SignatureMethod"] = "HmacSHA1"
	}
	sortedParamKeys := make([]string, 0)
	for k, _ := range param {
		sortedParamKeys = append(sortedParamKeys, k)
	}
	sort.Strings(sortedParamKeys)

	host := ""
	if strings.HasPrefix(this.Endpoint, "https") {
		host = this.Endpoint[8:]
	} else {
		host = this.Endpoint[7:]
	}

	src := this.Method + host + this.Path + "?"
	flag := false
	for _, key := range sortedParamKeys {
		if flag {
			src += "&"
		}
		src += key + "=" + param[key]
		flag = true
	}
	param["Signature"] = Sign(src, this.SecretKey, this.SignMethod)
	urlStr := ""
	reqStr := ""
	if this.Method == "GET" {
		urlStr = this.Endpoint + this.Path + "?"
		flag = false
		for k, v := range param {
			if flag {
				urlStr += "&"
			}
			urlStr += k + "=" + url.QueryEscape(v)
			flag = true
		}
		if len(urlStr) > 2048 {
			err = fmt.Errorf("url string length is large than 2048")
			return
		}
	} else {
		urlStr = this.Endpoint + this.Path
		flag := false
		for k, v := range param {
			if flag {
				reqStr += "&"
			}
			reqStr += k + "=" + url.QueryEscape(v)
			flag = true
		}
	}
	//log.Printf("urlStr :%v", urlStr)
	//log.Printf("reqStr :%v", reqStr)

	proxyUrlStr := this.Proxy
	if proxyUrl, found := param["proxyUrl"]; found {
		proxyUrlStr = proxyUrl
	}
	userTimeout := 0
	if UserpollingWaitSeconds, found := param["UserpollingWaitSeconds"]; found {
		userTimeout, err = strconv.Atoi(UserpollingWaitSeconds)
		if err != nil {
			return "", fmt.Errorf("strconv failed: %v", err.Error())
		}
	}

	resp, err = this.CmqHttp.request(this.Method, urlStr, reqStr, proxyUrlStr, userTimeout)
	if err != nil {
		return resp, fmt.Errorf("CmqHttp.request failed: %v", err.Error())
	}
	return
}
