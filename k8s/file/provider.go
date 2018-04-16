// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	teamvault_parser "github.com/bborbe/teamvault_utils/parser"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/seibert-media/k8s-deploy/k8s"
	k8s_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	k8s_scheme "k8s.io/client-go/kubernetes/scheme"
)

type walkFuncBuilder func([]k8s_runtime.Object) filepath.WalkFunc

type provider struct {
	templateDirectory TemplateDirectory
	parser            teamvault_parser.Parser
	walkFunc          walkFuncBuilder
}

// New file provider for directory using Teamvault parser
func New(
	templateDirectory TemplateDirectory,
	parser teamvault_parser.Parser,
) k8s.Provider {
	p := &provider{
		templateDirectory: templateDirectory,
		parser:            parser,
	}
	p.walkFunc = p.walkFuncBuilder
	return p
}

// GetObjects in the given namespace
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
	if err = filepath.Walk(dir.String(), p.walkFunc(result)); err != nil {
		return nil, fmt.Errorf("walk path failed: %v", err)
	}

	glog.V(1).Infof("read files completed. found %d objects", len(result))
	return result, nil

}

func (p *provider) walkFuncBuilder(result []k8s_runtime.Object) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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
		if obj, _, err = k8s_unstructured.UnstructuredJSONScheme.Decode(content, nil, obj); err != nil {
			return fmt.Errorf("unmarshal to object failed: %v", err)
		}

		glog.V(2).Infof("found file object %s", k8s.ObjectToString(obj))
		result = append(result, obj)
		glog.V(4).Infof("add object to result")
		return nil
	}
}

func kind(content []byte) (k8s_runtime.Object, error) {
	_, kind, err := k8s_unstructured.UnstructuredJSONScheme.Decode(content, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to unknown failed: %v", err)
	}
	obj, err := k8s_scheme.Scheme.New(*kind)
	if err != nil {
		return nil, fmt.Errorf("create object failed: %v", err)
	}
	return obj, nil
}
