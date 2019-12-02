package httpclient

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Requester struct {
	BaseUrl string
}

func (requester *Requester) RequestHttpByGet(endpoint string, params map[string]string) ([]byte, error){
	return requester.RequestHttp("GET",endpoint,params)
}

func (requester *Requester) RequestHttpByPost(endpoint string, params map[string]string) ([]byte, error){
	return requester.RequestHttp("POST",endpoint,params)
}


func (requester *Requester) RequestHttp(method string, endpoint string, params map[string]string) ([]byte, error) {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("%s/%s", requester.BaseUrl, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("build request error=%s",err.Error()))
	}
	req.WithContext(context.Background())

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}


	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("request alpha ping error")
	}

	return buf, nil
}
func (requester *Requester) PostString(endpoint string, reqstr string) ([]byte, error) {
	resp, err := http.Post(fmt.Sprintf("%s/rawtx", requester.BaseUrl,), "text/plain", strings.NewReader(reqstr))
	if err != nil {
		return nil,err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("request alpha ping error")
	}

	return buf, nil
}
