package downloader

import (
	"os"
)

type Downloader interface {
	Download(url string, targetDirectory *os.File) error
}
