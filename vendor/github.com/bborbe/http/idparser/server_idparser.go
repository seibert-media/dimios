package idparser

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func ParseIdFormRequest(request *http.Request) (int, error) {
	return ParseIdFromUri(request.RequestURI)
}

func ParseKeyFormRequest(request *http.Request) (string, error) {
	return ParseKeyFromUri(request.RequestURI)
}

func ParseIdFromUri(uri string) (int, error) {
	key, err := ParseKeyFromUri(uri)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(key)
}

func ParseKeyFromUri(uri string) (string, error) {
	re := regexp.MustCompile("[^?#]*/([^/\\?#]+?)([\\?#].*)?$")
	matches := re.FindStringSubmatch(uri)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("parse key from uri %s failed", uri)
}
