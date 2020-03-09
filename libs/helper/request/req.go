package CHelperRequest

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpCli struct {
}

var client = &http.Client{}

func NewHttpCli() interface{} { return new(HttpCli) }

func (this *HttpCli) DoForm(method string, url string, reqFunc func(r *http.Request), bodyValues url.Values) (http.Header, []byte, error) {
	var bodyReader io.Reader
	if bodyValues == nil {
		bodyReader = nil
	} else {
		bodyReader = strings.NewReader(bodyValues.Encode())
	}

	return this.Do(method, url, reqFunc, bodyReader)
}

func (this *HttpCli) DoJson(method string, url string, reqFunc func(r *http.Request), bodyValues map[string]string) (http.Header, []byte, error) {
	var bodyReader io.Reader
	if bodyValues == nil {
		bodyReader = nil
	} else {
		bytesData, err := json.Marshal(bodyValues)
		if err != nil {
			return nil, nil, err
		}
		bodyReader = strings.NewReader(string(bytesData))
	}

	return this.Do(method, url, reqFunc, bodyReader)
}

func (this *HttpCli) Do(method string, url string, reqFunc func(r *http.Request), bodyReader io.Reader) (http.Header, []byte, error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		return nil, nil, err
	}
	if reqFunc != nil {
		reqFunc(req)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return resp.Header, body, err
	}
}

func (this *HttpCli) JsonParse(json string) gjson.Result {
	return gjson.Parse(json)
}

var ErrNoRedirect = errors.New("Don't redirect!")
var clientNoRedirect = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return ErrNoRedirect
	},
}

func (this *HttpCli) DoOrigin(method string, url string, reqFunc func(r *http.Request), bodyReader io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		return nil, err
	}
	if reqFunc != nil {
		reqFunc(req)
	}

	return clientNoRedirect.Do(req)
}
