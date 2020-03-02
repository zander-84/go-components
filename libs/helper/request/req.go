package CHelperRequest

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpCli struct {
}

var client = &http.Client{}

func NewHttpCli() interface{} { return new(HttpCli) }

func (this *HttpCli) Do(method string, url string, reqFunc func(r *http.Request) *http.Request, bodyFunc func() io.Reader) ([]byte, error) {
	var body io.Reader
	if bodyFunc == nil {
		body = nil
	} else {
		body = bodyFunc()
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, body)
	if err != nil {
		return nil, err
	}
	if reqFunc != nil {
		tmpreq := reqFunc(req)
		if tmpreq != nil {
			req = tmpreq
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
}
