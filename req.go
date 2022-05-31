package frm

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

var(
	errIncorrectStatus = errors.New("resp status code != expected status code")
)

type FastHttp struct {
	Client         *fasthttp.Client
	Request        *fasthttp.Request
	Response       *fasthttp.Response
	Timeout        time.Duration
	ExpectedStatus int
	MaxRetries     int
	Proxy          *Proxy
}

type NetHttp struct {
	Client         *http.Client
	Request        *http.Request
	Response       *http.Response
	Timeout        time.Duration
	ExpectedStatus int
	MaxRetries     int
	Proxy          *Proxy
}

type FastData struct {
	Client         *fasthttp.Client
	Headers        map[string]string
	URL            string
	Method         string
	ExpectedStatus int
	Payload        string
	Timeout        int
	Proxy          *Proxy
}

type NetData struct {
	Client         *http.Client
	Headers        map[string]string
	URL            string
	Method         string
	ExpectedStatus int
	Payload        string
	Timeout        int
	Proxy          *Proxy
}

type Builder interface {
	Build()
}

type Requester interface {
	DoRequest()
}

// Send request via FastHttp
func (fh *FastHttp) DoRequest() (err error) {
	for i := 0; i < fh.MaxRetries + 1; i++ {
		if fh.Proxy.BansQuantity > fh.MaxRetries*10 {
			err = fmt.Errorf("too many bad requests with this proxy: %v", fh.Proxy.FastFmt)
			Logger.Error().Err(err)
			return err
		}
		if fh.Proxy != nil {
			fh.Client.Dial = fasthttpproxy.FasthttpHTTPDialer(fh.Proxy.FastFmt)
		}
		err = fh.Client.DoTimeout(fh.Request, fh.Response, fh.Timeout)
		if err != nil {
			Logger.Warn().Err(err).Str("req_data", fh.Request.String()).Str("proxy", fh.Proxy.FastFmt)
			continue
		}
		if fh.Response.StatusCode() != fh.ExpectedStatus {
			fh.Proxy.BansQuantity++
			Logger.Warn().Err(errIncorrectStatus).Str("req_data", fh.Request.String()).Str("resp_data", fh.Response.String()).Str("proxy", fh.Proxy.FastFmt)
			continue
		}
		return nil
	}
	// Logger.Error().Err(errors.New("couldn't do request")).Str("req_data", fh.Request.String()).Str("resp_data", fh.Response.String())
	return errIncorrectStatus
}

// // Send request via net/http
// func (nh *NetHttp) DoRequest() (err error) {
// 	for i := 0; i < nh.MaxRetries; i++ {
// 		nh.Response, err = nh.Client.Do(nh.Request)
// 		if err != nil {
// 			Logger.Warn().Err(err).Str("data", nh.Request.String())
// 			continue
// 		}
// 		defer nh.Response.Body.Close()
// 	}
// 	return nil
// }

//Build data for FastHttp
func (data *FastData) Build() *FastHttp {
	if data.Client == nil {
		data.Client = &fasthttp.Client{
			MaxIdleConnDuration:           5 * time.Second,
			ReadBufferSize:                8192,
			DisableHeaderNamesNormalizing: true,
		}
	}

	req := fasthttp.AcquireRequest()
	req.Header.SetMethodBytes([]byte(data.Method))
	req.SetRequestURIBytes([]byte(data.URL))

	switch data.Method {
	case "POST", "PUT":
		req.SetBody([]byte(data.Payload))
	}

	for k, v := range data.Headers {
		req.Header.SetBytesKV([]byte(k), []byte(v))
	}
	res := fasthttp.AcquireResponse()

	if data.Timeout == 0 {
		data.Timeout = 3
	}

	return &FastHttp{Client: data.Client, Request: req, Response: res, Timeout: time.Duration(data.Timeout) * time.Second, ExpectedStatus: data.ExpectedStatus}
}

// //Build data for net/http
// func (data *NetData) Build() *NetHttp {
// 	transport := &http.Transport{
// 		TLSClientConfig: &tls.Config{
// 			MinVersion:               tls.VersionTLS12,
// 			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
// 			PreferServerCipherSuites: true,
// 			CipherSuites: []uint16{
// 				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
// 				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
// 				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
// 				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
// 				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
// 			},
// 		},
// 	}
// 	if data.Timeout == 0 {
// 		data.Timeout = 3
// 	}
//
// 	client := &http.Client{
// 		Transport: transport,
// 		Timeout:   time.Duration(data.Timeout) * time.Second,
// 	}
//
// 	switch data.Method {
// 	case "POST", "PUT":
// 		payload := strings.NewReader(data.Payload)
// 		req, err := http.NewRequest(data.Method, data.URL, payload)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
//
// 		for k, v := range data.Headers {
// 			req.Header.Set(k, v)
// 		}
//
// 		return &NetHttp{Client: client, Request: req, Response: &http.Response{}, Timeout: time.Duration(data.Timeout) * time.Second}
// 	}
//
// 	req, err := http.NewRequest(data.Method, data.URL, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	for k, v := range data.Headers {
// 		req.Header.Set(k, v)
// 	}
//
// 	return &NetHttp{Client: client, Request: req, Response: &http.Response{}, Timeout: time.Duration(data.Timeout) * time.Second}
// }
