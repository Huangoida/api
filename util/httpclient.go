package util

import (
	"github.com/valyala/fasthttp"
	"time"
)

func getHttpClient() *fasthttp.Client {

	reqClient := &fasthttp.Client{
		ReadTimeout:                   time.Second * 5,
		WriteTimeout:                  time.Second * 5,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			// 最大并发数，0表示无限制
			Concurrency: 4096,
			// 将 DNS 缓存时间从默认分钟增加到一小时
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	return reqClient
}

func Do(method string, url string, headers map[string]string, parameters map[string]string, body []byte, contentType string) (string, int) {
	client := getHttpClient()
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.Header.SetMethod(getMethod(method))
	req.SetRequestURI(url)
	var arg fasthttp.Args
	for k, v := range parameters {
		arg.Add(k, v)
	}
	req.URI().SetQueryString(arg.String())
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if len(body) != 0 {
		req.SetBodyRaw([]byte(body))
	}
	if contentType != "" {
		req.Header.SetContentType(contentType)
	}

	if err := client.Do(req, resp); err != nil {
		return err.Error(), resp.StatusCode()
	}
	return string(resp.Body()), resp.StatusCode()
}

func getMethod(method string) string {
	switch method {
	case "GET":
		return fasthttp.MethodGet
	case "POST":
		return fasthttp.MethodPost
	case "PUT":
		return fasthttp.MethodPut
	case "DELETE":
		return fasthttp.MethodDelete
	default:
		return method
	}
}
