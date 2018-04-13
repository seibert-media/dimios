package redirect_follower

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang/glog"
)

const LIMIT = 10

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type RedirectFollower interface {
	ExecuteRequestAndFollow(req *http.Request) (resp *http.Response, err error)
}

type redirectFollower struct {
	executeRequest ExecuteRequest
}

func New(executeRequest ExecuteRequest) *redirectFollower {
	r := new(redirectFollower)
	r.executeRequest = executeRequest
	return r
}

func (r *redirectFollower) ExecuteRequestAndFollow(req *http.Request) (*http.Response, error) {
	return executeRequestAndFollow(r.executeRequest, req, 0)
}

func executeRequestAndFollow(executeRequest ExecuteRequest, req *http.Request, counter int) (*http.Response, error) {
	glog.V(4).Infof("execute request to %s", req.URL)
	glog.V(4).Infof("request %v\n", req)
	resp, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("response %v", resp)
	if resp.StatusCode/100 == 3 {
		glog.V(4).Infof("redirect - statuscode: %d", resp.StatusCode)
		if counter > LIMIT {
			return nil, fmt.Errorf("redirect limit reached")
		}
		var reqCopy http.Request = *req
		var p *http.Request = &reqCopy
		var location []string = resp.Header["Location"]
		if len(location) != 1 {
			return nil, fmt.Errorf("redirect failed")
		}
		glog.V(4).Infof("redirect to %s", location[0])
		p.URL, err = locationToUrl(req.URL, location[0])
		if err != nil {
			return nil, err
		}
		p.Host = p.URL.Host
		return executeRequestAndFollow(executeRequest, p, counter+1)
	}

	return resp, nil
}

func locationToUrl(u *url.URL, location string) (*url.URL, error) {
	if len(location) == 0 {
		return nil, fmt.Errorf("empty location")
	}
	if location[0] == '/' {
		location = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, location)
	}
	return url.Parse(location)
}
