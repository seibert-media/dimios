package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bborbe/teamvault_utils/model"
	"github.com/bborbe/teamvault_utils/parser"
	"github.com/golang/glog"
)

type configGenerator struct {
	configParser parser.Parser
}

func New(
	configParser parser.Parser,
) *configGenerator {
	c := new(configGenerator)
	c.configParser = configParser
	return c
}

func (c *configGenerator) Generate(sourceDirectory model.SourceDirectory, targetDirectory model.TargetDirectory) error {
	glog.V(4).Infof("generate config from %s to %s", sourceDirectory.String(), targetDirectory.String())
	return filepath.Walk(sourceDirectory.String(), func(path string, info os.FileInfo, err error) error {
		glog.V(4).Infof("generate path %s info %v", path, info)
		if err != nil {
			return err
		}
		target := fmt.Sprintf("%s%s", targetDirectory.String(), strings.TrimPrefix(path, sourceDirectory.String()))
		glog.V(2).Infof("target: %s", target)
		if info.IsDir() {
			err := os.MkdirAll(target, 0755)
			if err != nil {
				glog.V(2).Infof("create directory %s failed: %v", target, err)
				return err
			}
			glog.V(4).Infof("directory %s created", target)
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			glog.V(2).Infof("read file %s failed: %v", path, err)
			return err
		}
		content, err = c.configParser.Parse(content)
		if err != nil {
			glog.V(2).Infof("replace variables failed: %v", err)
			return err
		}
		if err := ioutil.WriteFile(target, content, 0644); err != nil {
			glog.V(2).Infof("create file %s failed: %v", target, err)
			return err
		}
		glog.V(4).Infof("file %s created", target)
		return nil
	})
}
