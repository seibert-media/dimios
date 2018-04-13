package util

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bborbe/io/reader_shadow_copy"
	"github.com/golang/glog"
)

func ResponseToString(response *http.Response) (string, error) {
	content, err := ResponseToByteArray(response)
	return string(content), err
}

func ResponseToByteArray(response *http.Response) ([]byte, error) {
	body := response.Body
	return ioutil.ReadAll(body)
}

var contentTypeToExt = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
}

func FindFileExtension(response *http.Response) (string, error) {
	var ext string
	if response.Request != nil && response.Request.URL != nil {
		path := response.Request.URL.Path
		pos := strings.LastIndex(path, ".")
		if pos > 0 {
			ext = path[pos+1:]
		}
	}
	if response.Header != nil {
		contentType := response.Header.Get("Content-Type")
		ext = contentTypeToExt[contentType]
	}
	if len(ext) == 0 {
		return "", fmt.Errorf("find extension failed")
	}
	return ext, nil
}

// PrintDump prints dump of request, optionally writing it in the response
func PrintDump(request *http.Request) {
	glog.V(4).Infof("header: %v", request.Header)
	dump, _ := httputil.DumpRequest(request, true)
	glog.V(4).Infof("request: %v", string(dump))
}

// Decode into a ma[string]interface{} the JSON in the POST Request
func DecodePostJSON(request *http.Request, logging bool) (map[string]interface{}, error) {
	var err error
	var payLoad map[string]interface{}
	if logging {
		reader := reader_shadow_copy.New(request.Body)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&payLoad)
		glog.V(4).Infof("body: %s", string(reader.Bytes()))
		return payLoad, err
	}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&payLoad)
	return payLoad, err
}
