package main

import (
	"bufio"
	"os"

	"flag"
	"io"

	"github.com/golang/glog"

	"fmt"
	"runtime"
	"sync"

	http_client_builder "github.com/bborbe/http/client_builder"
	http_downloader "github.com/bborbe/http/downloader"
	http_downloader_by_url "github.com/bborbe/http/downloader/by_url"
	io_util "github.com/bborbe/io/util"
)

const (
	PARAMETER_PARALLEL_DOWNLOADS = "max"
	PARAMETER_TARGET             = "target"
	DEFAULT_PARALLEL_DOWNLOADS   = 2
	DEFAULT_TARGET               = "~/Downloads"
)

var (
	maxConcurrencyDownloadsPtr = flag.Int(PARAMETER_PARALLEL_DOWNLOADS, DEFAULT_PARALLEL_DOWNLOADS, "max parallel downloads")
	targetDirectoryPtr         = flag.String(PARAMETER_TARGET, DEFAULT_TARGET, "directory")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	input := os.Stdin
	wg := new(sync.WaitGroup)
	httpClientBuilder := http_client_builder.New()
	httpClient := httpClientBuilder.Build()
	downloader := http_downloader_by_url.New(httpClient.Get)

	err := do(writer, input, *maxConcurrencyDownloadsPtr, wg, downloader, *targetDirectoryPtr)
	wg.Wait()
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer, input io.Reader, maxConcurrencyDownloads int, wg *sync.WaitGroup, downloader http_downloader.Downloader, targetDirectoryName string) error {
	var err error
	if targetDirectoryName, err = io_util.NormalizePath(targetDirectoryName); err != nil {
		return err
	}
	if isDir, err := io_util.IsDirectory(targetDirectoryName); err != nil || isDir == false {
		fmt.Fprintf(writer, "parameter %s is invalid\n", PARAMETER_TARGET)
		return fmt.Errorf("parameter is not a directory")
	}
	targetDirectory, err := os.Open(targetDirectoryName)
	if err != nil {
		return err
	}
	throttle := make(chan bool, maxConcurrencyDownloads)
	reader := bufio.NewReader(input)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			link := string(line)
			throttle <- true
			if err := downloader.Download(link, targetDirectory); err != nil {
				glog.Warningf("download failed: %v", err)
			}
			<-throttle
			wg.Done()
		}()
	}
}
