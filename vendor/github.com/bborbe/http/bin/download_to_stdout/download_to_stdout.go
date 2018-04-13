package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"flag"

	"runtime"

	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	input := os.Stdin
	err := do(writer, input)
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer, input io.Reader) error {
	reader := bufio.NewReader(input)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = downloadLink(writer, string(line))
		if err != nil {
			return err
		}
	}
}

func downloadLink(writer io.Writer, url string) error {
	glog.V(4).Infof("download %s started", url)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if glog.V(4) {
			glog.Infof("%s", string(content))
		}
		return errors.New(string(content))
	}
	if _, err := io.Copy(writer, response.Body); err != nil {
		return err
	}
	glog.V(4).Infof("download %s finished", url)
	return nil
}
