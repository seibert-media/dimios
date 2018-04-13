package md5

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/http/downloader"
)

func TestImplementsDownloader(t *testing.T) {
	r := New(nil)
	var i *downloader.Downloader
	if err := AssertThat(r, Implements(i)); err != nil {
		t.Fatal(err)
	}
}
