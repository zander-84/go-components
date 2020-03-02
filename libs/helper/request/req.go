package CHelperRequest

import (
	"encoding/json"
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

func (this *HttpCli) DoForm(method string, url string, reqFunc func(r *http.Request), bodyValues url.Values) ([]byte, error) {
	var bodyReader io.Reader
	if bodyValues == nil {
		bodyReader = nil
	} else {
		bodyReader = strings.NewReader(bodyValues.Encode())
	}

	return this.do(method, url, reqFunc, bodyReader)
}

func (this *HttpCli) DoJson(method string, url string, reqFunc func(r *http.Request), bodyValues map[string]string) ([]byte, error) {
	var bodyReader io.Reader
	if bodyValues == nil {
		bodyReader = nil
	} else {
		bytesData, err := json.Marshal(bodyValues)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(string(bytesData))
	}

	return this.do(method, url, reqFunc, bodyReader)
}

func (this *HttpCli) do(method string, url string, reqFunc func(r *http.Request), bodyReader io.Reader) ([]byte, error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		return nil, err
	}
	if reqFunc != nil {
		reqFunc(req)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
}

func (this *HttpCli) JsonParse(json string) gjson.Result {
	return gjson.Parse(json)
}
