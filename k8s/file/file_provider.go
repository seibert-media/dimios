// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	teamvault_parser "github.com/bborbe/teamvault-utils/parser"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/seibert-media/dimios/k8s"
	k8s_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	k8s_scheme "k8s.io/client-go/kubernetes/scheme"
)

type provider struct {
	templateDirectory TemplateDirectory
	parser            teamvault_parser.Parser
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
	err = filepath.Walk(dir.String(), func(path string, info os.FileInfo, err error) error {
		glog.V(4).Infof("handle path: %s", path)

		if info.IsDir() {
			glog.V(3).Infof("skip directory: %s", path)
			return nil
		}
		yamlTemplate, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s failed: %v", path, err)
		}
		if glog.V(6) {
			glog.Infof("yaml: %s", yamlTemplate)
		}

		yamlContent, err := p.parser.Parse(yamlTemplate)
		if err != nil {
			if glog.V(4) {
				glog.Infof("content: %s", string(yamlTemplate))
			}
			return fmt.Errorf("parse content of file %s failed: %v", path, err)
		}
		if glog.V(6) {
			glog.Infof("parsed: %s", yamlContent)
		}

		jsonContent, err := yaml.YAMLToJSON(yamlContent)
		if err != nil {
			if glog.V(4) {
				glog.Infof("content: %s", string(yamlContent))
			}
			return fmt.Errorf("convert yaml to json for file %s failed: %v", path, err)
		}
		if glog.V(6) {
			glog.Infof("json: %s", jsonContent)
		}

		obj, err := kind(jsonContent)
		if err != nil {
			if glog.V(4) {
				glog.Infof("content: %s", string(jsonContent))
			}
			return fmt.Errorf("create object by content for file %s failed: %v", path, err)
		}

		glog.V(4).Infof("found kind %v", obj.GetObjectKind())
		if obj, _, err = k8s_unstructured.UnstructuredJSONScheme.Decode(jsonContent, nil, obj); err != nil {
			return fmt.Errorf("unmarshal to object for file %s failed: %v", path, err)
		}

		glog.V(2).Infof("found object %s in file %s", k8s.ObjectToString(obj), path)
		result = append(result, obj)
		glog.V(4).Infof("add object to result for file %s", path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk path failed: %v", err)
	}
	glog.V(1).Infof("found in files %d object for namespace %s", len(result), namespace)
	return result, nil
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
