package header

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"
)

func CreateAuthorizationToken(name string, value string) string {
	glog.V(4).Infof("create bearer")
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", name, value)))
}

func CreateAuthorizationTokenSimple(name string) string {
	glog.V(4).Infof("create simple bearer")
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func ParseAuthorizationToken(token string) (string, string, error) {
	glog.V(4).Infof("parse token")
	value, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", err
	}
	parse := strings.SplitN(string(value), ":", 2)
	if len(parse) != 2 {
		return "", "", fmt.Errorf("parse header failed")
	}
	return parse[0], parse[1], nil
}

func ParseAuthorizationTokenSimple(token string) (string, error) {
	glog.V(4).Infof("parse token")
	value, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func CreateAuthorizationBearerHeader(name string, value string) string {
	return fmt.Sprintf("Bearer %s", CreateAuthorizationToken(name, value))
}

func CreateAuthorizationBearerHeaderSimple(name string) string {
	return fmt.Sprintf("Bearer %s", CreateAuthorizationTokenSimple(name))
}

func ParseAuthorizationBearerHttpRequest(req *http.Request) (string, string, error) {
	return ParseAuthorizationHttpRequest("Bearer", req)
}

func ParseAuthorizationBearerHttpRequestSimple(req *http.Request) (string, error) {
	return ParseAuthorizationHttpRequestSimple("Bearer", req)
}

func ParseAuthorizationBasisHttpRequest(req *http.Request) (string, string, error) {
	return ParseAuthorizationHttpRequest("Basic", req)
}

func ParseAuthorizationHttpRequest(authtype string, req *http.Request) (string, string, error) {
	authorizations := req.Header["Authorization"]
	if len(authorizations) != 1 {
		return "", "", fmt.Errorf("header Authorization invalid")
	}
	return ParseAuthorizationHeader(authtype, authorizations[0])
}

func ParseAuthorizationHttpRequestSimple(authtype string, req *http.Request) (string, error) {
	authorizations := req.Header["Authorization"]
	if len(authorizations) != 1 {
		return "", fmt.Errorf("header Authorization invalid")
	}
	return ParseAuthorizationHeaderSimple(authtype, authorizations[0])
}

func ParseAuthorizationHeader(authtype string, header string) (string, string, error) {
	glog.V(4).Infof("parse %s", authtype)
	if strings.Index(header, fmt.Sprintf("%s ", authtype)) != 0 {
		return "", "", fmt.Errorf("header Authorization invalid")
	}
	tokens := strings.SplitN(header, " ", 2)
	if len(tokens) != 2 {
		return "", "", fmt.Errorf("parse token from header failed")
	}
	return ParseAuthorizationToken(tokens[1])
}

func ParseAuthorizationHeaderSimple(authtype string, header string) (string, error) {
	glog.V(4).Infof("parse %s", authtype)
	if strings.Index(header, fmt.Sprintf("%s ", authtype)) != 0 {
		return "", fmt.Errorf("header Authorization invalid")
	}
	tokens := strings.SplitN(header, " ", 2)
	if len(tokens) != 2 {
		return "", fmt.Errorf("parse token from header failed")
	}
	return ParseAuthorizationTokenSimple(tokens[1])
}
