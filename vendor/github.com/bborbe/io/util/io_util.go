package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

func IsDirectory(dir string) (bool, error) {
	glog.V(4).Infof("IsDir %s", dir)
	file, err := os.Open(dir)
	defer file.Close()
	if err != nil {
		glog.V(2).Infof("IsDir - open dir %s failed: %v", dir, err)
		return false, nil
	}
	fileinfo, err := file.Stat()
	if err != nil {
		glog.V(2).Infof("IsDir get state for dir %s failed: %v", dir, err)
		return false, err
	}
	return fileinfo.IsDir(), nil
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func NormalizePath(path string) (string, error) {
	glog.V(4).Infof("NormalizePath %s", path)
	if strings.Index(path, "~/") == 0 {
		home := os.Getenv("HOME")
		if len(home) == 0 {
			glog.V(2).Infof("normalize path failed, enviroment variable HOME missing")
			return "", fmt.Errorf("env HOME not found")
		}
		path = fmt.Sprintf("%s/%s", home, path[2:])
		glog.V(2).Infof("replace ~/ with homedir. new path: %s", path)
	}
	result, err := filepath.Abs(path)
	if err != nil {
		glog.Warningf("get absolute path for %v failed: %v", path, err)
		return "", err
	}
	return result, nil
}
