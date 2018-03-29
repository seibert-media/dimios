package find

import (
	"context"
	"os"
	"path/filepath"
	"github.com/bborbe/k8s_deploy/change"
	"github.com/golang/glog"
	"fmt"
	io_util "github.com/bborbe/io/util"
)

type ManifestDirectory string

func (m ManifestDirectory) String() string {
	return string(m)
}

type Finder interface {
	Changes(ctx context.Context, c chan<- change.Change) error
}

type finder struct {
	dir ManifestDirectory
}

func NewFinder(dir ManifestDirectory) Finder {
	return &finder{
		dir: dir,
	}
}

func (f *finder) Changes(ctx context.Context, c chan<- change.Change) error {
	defer close(c)
	path, err := io_util.NormalizePath(f.dir.String())
	if err != nil {
		return fmt.Errorf("create abs path of %s failed: %v", path, err)
	}
	return filepath.Walk(
		path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				glog.V(4).Infof("walk %s failed: %v", path, err)
				return err
			}
			if info.IsDir() {
				glog.V(4).Infof("%s is a dir => skip", path)
				return nil
			}
			if filepath.Ext(path) != ".yaml" {
				glog.V(4).Infof("%s is a not a yaml => skip", path)
				return nil
			}
			glog.V(4).Infof("send to %s as change", path)
			return sendChange(ctx, c, NewFileChange())
		},
	)
}

func sendChange(ctx context.Context, c chan<- change.Change, change change.Change) error {
	select {
	case c <- change:
		glog.V(4).Infof("added %v to channel", change)
		return nil
	case <-ctx.Done():
		glog.V(3).Infoln("context done, skip add changes")
		return nil
	}
}

func NewFileChange() change.Change {
	return &struct{}{}
}
