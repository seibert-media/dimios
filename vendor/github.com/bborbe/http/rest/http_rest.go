package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"io"

	"net/url"

	"github.com/golang/glog"
)

type executeRequest func(req *http.Request) (resp *http.Response, err error)

type rest struct {
	executeRequest executeRequest
}

type Rest interface {
	Call(url string, values url.Values, method string, request interface{}, response interface{}, headers http.Header) error
}

func New(
	executeRequest executeRequest,
) *rest {
	r := new(rest)
	r.executeRequest = executeRequest
	return r
}

func (r *rest) Call(url string, values url.Values, method string, request interface{}, response interface{}, headers http.Header) error {
	if values != nil {
		url = fmt.Sprintf("%s?%s", url, values.Encode())
	}
	glog.V(4).Infof("call %s on path %s", method, url)
	start := time.Now()
	defer glog.V(4).Infof("create completed in %dms", time.Now().Sub(start)/time.Millisecond)
	glog.V(4).Infof("send message to %s", url)

	var body io.Reader
	if request != nil {
		content, err := json.Marshal(request)
		if err != nil {
			glog.V(2).Infof("marhal request failed: %v", err)
			return err
		}
		if glog.V(4) {
			glog.Infof("send request to %s: %s", url, string(content))
		}
		body = bytes.NewBuffer(content)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		glog.V(2).Infof("build request failed: %v", err)
		return err
	}
	req.Header.Set("ContentType", "application/json")
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	resp, err := r.executeRequest(req)
	if err != nil {
		glog.V(2).Infof("execute request failed: %v", err)
		return err
	}
	if resp.StatusCode/100 != 2 {
		glog.V(2).Infof("status %d not 2xx", resp.StatusCode)
		return fmt.Errorf("request to %s failed with status: %d", url, resp.StatusCode)
	}
	if response != nil {
		if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
			glog.V(2).Infof("decode response failed: %v", err)
			return err
		}
	}
	glog.V(4).Infof("rest call successful")
	return nil
}
