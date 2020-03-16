package CHelperRequest

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var Attempt = 2

type HttpCli struct {
	Attempt int
}

func NewHttpCli() interface{} {
	this := new(HttpCli)
	this.Attempt = Attempt
	return this
}

func (this *HttpCli) DoValues(method string, url string, reqFunc func(r *http.Request), bodyValues url.Values) (http.Header, []byte, error) {
	return this.doValues(method, url, reqFunc, bodyValues)
}

func (this *HttpCli) doValues(method string, url string, reqFunc func(r *http.Request), bodyValues url.Values) (http.Header, []byte, error) {
	var bodyReader io.Reader
	if bodyValues == nil {
		bodyReader = nil
	} else {
		bodyReader = strings.NewReader(bodyValues.Encode())
	}

	return this.Do(method, url, reqFunc, bodyReader)
}

func (this *HttpCli) DoMap(method string, url string, reqFunc func(r *http.Request), bodyValues map[string]string) (http.Header, []byte, error) {
	return this.doMap(method, url, reqFunc, bodyValues)
}

func (this *HttpCli) doMap(method string, url string, reqFunc func(r *http.Request), bodyValues map[string]string) (http.Header, []byte, error) {
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
	var client = &http.Client{}

	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		return nil, nil, err
	}
	if reqFunc != nil {
		reqFunc(req)
	}
	req.Close = true

	//It looks like the that server (Apache 1.3, wow!) is serving up a truncated gzip response. If you explicitly request the identity encoding (preventing the Go transport from adding gzip itself), you won't get the ErrUnexpectedEOF:
	//req.Header.Add("Accept-Encoding", "identity")

	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		return resp.Header, body, err
	}
}

var ErrNoRedirect = errors.New("Don't redirect!")

func (this *HttpCli) DoOrigin(method string, url string, reqFunc func(r *http.Request), bodyReader io.Reader) (*http.Response, error) {
	var clientNoRedirect = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrNoRedirect
		},
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		return nil, err
	}
	if reqFunc != nil {
		reqFunc(req)
	}

	return clientNoRedirect.Do(req)
}

func (this *HttpCli) JsonParse(json string) gjson.Result {
	return gjson.Parse(json)
}

func (this *HttpCli) GetUrlForHtmlParse(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return goquery.NewDocumentFromReader(res.Body)
}
