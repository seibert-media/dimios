package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/k8s_deploy/k8s"
	"github.com/bborbe/teamvault_utils/parser"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type TemplateDirectory string

func (t TemplateDirectory) String() string {
	return string(t)
}

func (d TemplateDirectory) NormalizePath() (TemplateDirectory, error) {
	root, err := io_util.NormalizePath(d.String())
	if err != nil {
		return "", err
	}
	return TemplateDirectory(root), nil
}

func (t *TemplateDirectory) PathToNamespace(namespace k8s.Namespace) NamespaceDirectory {
	return NamespaceDirectory(path.Join(t.String(), namespace.String()))
}

type NamespaceDirectory string

func (n NamespaceDirectory) String() string {
	return string(n)
}

func (n NamespaceDirectory) Exists() bool {
	f, err := os.Open(n.String())
	if err != nil {
		return false
	}
	fs, err := f.Stat()
	if err != nil {
		return false
	}
	return fs.IsDir()
}

type provider struct {
	templateDirectory TemplateDirectory
	parser            parser.Parser
}

func New(
	templateDirectory TemplateDirectory,
	parser parser.Parser,
) k8s.Provider {
	return &provider{
		templateDirectory: templateDirectory,
		parser:            parser,
	}
}

func (p *provider) GetObjects(namespace k8s.Namespace) ([]k8s_runtime.Object, error) {
	path, err := p.templateDirectory.NormalizePath()
	if err != nil {
		return nil, fmt.Errorf("normalize template directory failed: %v", err)
	}
	dir := path.PathToNamespace(namespace)
	if !dir.Exists() {
		return nil, fmt.Errorf("namespace %s not found", namespace)
	}
	var result []k8s_runtime.Object
	err = filepath.Walk(dir.String(), func(path string, info os.FileInfo, err error) error {
		glog.V(4).Infof("handle path: %s", path)
		if info.IsDir() {
			glog.V(3).Infof("skip directory: %s", path)
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s failed: %v", path, err)
		}
		content, err = p.parser.Parse(content)
		if err != nil {
			return fmt.Errorf("parse content failed: %v", err)
		}
		if glog.V(6) {
			glog.Infof("yaml %s", content)
		}
		content, err = yaml.YAMLToJSON(content)
		if err != nil {
			return fmt.Errorf("yaml to json failed: %v", err)
		}
		if glog.V(6) {
			glog.Infof("json %s", content)
		}
		glog.V(4).Infof("parse teamvault secrets to content completed")
		obj, err := kind(content)
		if err != nil {
			return fmt.Errorf("create object by content failed: %v", err)
		}
		glog.V(4).Infof("found kind %v", obj.GetObjectKind())
		if obj, _, err = unstructured.UnstructuredJSONScheme.Decode(content, nil, obj); err != nil {
			return fmt.Errorf("unmarshal to object failed: %v", err)
		}
		glog.V(2).Infof("found file object %s", k8s.ObjectToString(obj))
		result = append(result, obj)
		glog.V(4).Infof("add object to result")
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk path failed: %v", err)
	}
	glog.V(1).Infof("read files completed. found %d objects", len(result))
	return result, nil

}

func kind(content []byte) (k8s_runtime.Object, error) {
	_, kind, err := unstructured.UnstructuredJSONScheme.Decode(content, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to unknown failed: %v", err)
	}
	obj, err := scheme.Scheme.New(*kind)
	if err != nil {
		return nil, fmt.Errorf("create object failed: %v", err)
	}
	return obj, nil
}
