package wx

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type HttpClient struct {
	client    *http.Client
	userAgent string
}

func NewHttpClient() *HttpClient {
	var netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial:  (&net.Dialer{Timeout: 100 * time.Second}).Dial,
		// TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		// TLSHandshakeTimeout: 100 * time.Second,
	}
	cookieJar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Timeout:   time.Second * 100,
		Transport: netTransport,
		Jar:       cookieJar,
	}

	return &HttpClient{
		client:    httpClient,
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36 ",
	}
}

func (c *HttpClient) SetJar(jar http.CookieJar) {
	c.client.Jar = jar
}

func (c *HttpClient) Get(url string, data *url.Values) ([]byte, error) {
	if data != nil {
		url = url + "?" + data.Encode()
	}
	return c.fetch("GET", url, []byte(""), Header{})
}

func (c *HttpClient) GetByte(url string, data []byte) ([]byte, error) {

	return c.fetch("GET", url, data, Header{})
}

func (c *HttpClient) GetWithHeader(url string, heder Header) ([]byte, error) {

	return c.fetch("GET", url, nil, heder)
}

func (c *HttpClient) Post(url string, data *url.Values) ([]byte, error) {
	return c.fetch("POST", url, []byte(data.Encode()), Header{"Content-Type": "application/x-www-form-urlencoded"})
}

func (c *HttpClient) PostJson(url string, m map[string]interface{}) ([]byte, error) {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return c.fetch("POST", url, jsonString, Header{"Content-Type": "application/json; charset=UTF-8"})
}

func (c *HttpClient) PostJsonByte(url string, json []byte) ([]byte, error) {
	return c.fetch("POST", url, json, Header{"Content-Type": "application/json; charset=UTF-8"})
}

func (c *HttpClient) PostJsonByteForResp(url string, json []byte) (*http.Response, []byte, error) {
	return c.fetchResp("POST", url, json, Header{"Content-Type": "application/json; charset=UTF-8"})
}

func (c *HttpClient) FetchReponse(method string, uri string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.client.Do(req)
}

func (c *HttpClient) fetchReponseWithReader(method string, uri string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.client.Do(req)
}

func (c *HttpClient) fetchWithReader(method string, uri string, body io.Reader, headers Header) ([]byte, error) {
	resp, err := c.fetchReponseWithReader(method, uri, body, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *HttpClient) fetch(method string, uri string, body []byte, headers Header) ([]byte, error) {
	resp, err := c.FetchReponse(method, uri, body, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *HttpClient) fetchResp(method string, uri string, body []byte, headers Header) (resp *http.Response, b []byte, err error) {
	resp, err = c.FetchReponse(method, uri, body, headers)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	b, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, nil, err2
	}
	return resp, b, nil
}
