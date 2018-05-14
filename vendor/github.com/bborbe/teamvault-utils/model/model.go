package model

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	"fmt"
	"strings"

	io_util "github.com/bborbe/io/util"
	"github.com/golang/glog"
)

type VariableName string

func (v VariableName) String() string {
	return string(v)
}

type TeamvaultKey string

func (t TeamvaultKey) String() string {
	return string(t)
}

type SourceDirectory string

func (s SourceDirectory) String() string {
	return string(s)
}

type TargetDirectory string

func (t TargetDirectory) String() string {
	return string(t)
}

type Staging bool

func (s Staging) Bool() bool {
	return bool(s)
}

type TeamvaultUrl string

func (t TeamvaultUrl) String() string {
	return string(t)
}

type TeamvaultUser string

func (t TeamvaultUser) String() string {
	return string(t)
}

type TeamvaultPassword string

func (t TeamvaultPassword) String() string {
	return string(t)
}

type TeamvaultCurrentRevision string

func (t TeamvaultCurrentRevision) String() string {
	return string(t)
}

type TeamvaultFile string

func (t TeamvaultFile) String() string {
	return string(t)
}

func (t TeamvaultFile) Content() ([]byte, error) {
	return base64.StdEncoding.DecodeString(t.String())
}

type TeamvaultConfig struct {
	Url      TeamvaultUrl      `json:"url"`
	User     TeamvaultUser     `json:"user"`
	Password TeamvaultPassword `json:"pass"`
}

type TeamvaultConfigPath string

func (t TeamvaultConfigPath) String() string {
	return string(t)
}

func (d TeamvaultConfigPath) NormalizePath() (TeamvaultConfigPath, error) {
	root, err := io_util.NormalizePath(d.String())
	if err != nil {
		return "", err
	}
	return TeamvaultConfigPath(root), nil
}

// Exists the backup
func (t TeamvaultConfigPath) Exists() bool {
	path, err := t.NormalizePath()
	if err != nil {
		glog.V(2).Infof("normalize path failed: %v", err)
		return false
	}
	fileInfo, err := os.Stat(path.String())
	if err != nil {
		glog.V(2).Infof("file %v exists => false", t)
		return false
	}
	if fileInfo.Size() == 0 {
		glog.V(2).Infof("file %v empty => false", t)
		return false
	}
	if fileInfo.IsDir() {
		glog.V(2).Infof("file %v is dir => false", t)
		return false
	}
	glog.V(2).Infof("file %v exists and not empty => true", t)
	return true
}

func (t TeamvaultConfigPath) Parse() (*TeamvaultConfig, error) {
	path, err := t.NormalizePath()
	if err != nil {
		glog.V(2).Infof("normalize path failed: %v", err)
		return nil, err
	}
	content, err := ioutil.ReadFile(path.String())
	if err != nil {
		glog.Warningf("read config from file %v failed: %v", t, err)
		return nil, err
	}
	return ParseTeamvaultConfig(content)
}

func ParseTeamvaultConfig(content []byte) (*TeamvaultConfig, error) {
	config := &TeamvaultConfig{}
	if err := json.Unmarshal(content, config); err != nil {
		glog.Warningf("parse config failed: %v", err)
		return nil, err
	}
	return config, nil
}

type TeamvaultApiUrl string

func (t TeamvaultApiUrl) String() string {
	return string(t)
}

func (t TeamvaultApiUrl) Key() (TeamvaultKey, error) {
	parts := strings.Split(t.String(), "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("parse key form api-url failed")
	}
	return TeamvaultKey(parts[len(parts)-2]), nil
}
