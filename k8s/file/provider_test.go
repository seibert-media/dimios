// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/teamvault_utils/parser"
	"github.com/seibert-media/k8s-deploy/k8s"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

func TestGetObjectsNamespaceErrors(t *testing.T) {
	p := New("/tmp", nil)
	_, err := p.GetObjects("invalid-ns")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err.Error(), Is("namespace invalid-ns not found")); err != nil {
		t.Fatal(err)
	}
}

func Test_provider_GetObjects(t *testing.T) {
	type fields struct {
		templateDirectory TemplateDirectory
		parser            parser.Parser
		walkFunc          walkFuncBuilder
	}
	tests := []struct {
		name      string
		fields    fields
		namespace k8s.Namespace
		want      []k8s_runtime.Object
		wantErr   bool
	}{
		{
			"invalid-ns",
			fields{"/tmp", nil, nil},
			"invalid-ns",
			nil,
			true,
		},
		{
			"walk",
			fields{"/", nil, func(result []k8s_runtime.Object) filepath.WalkFunc {
				return func(path string, info os.FileInfo, err error) error {
					return nil
				}
			}},
			"tmp",
			nil,
			false,
		},
		{
			"walk-error",
			fields{"/", nil, func(result []k8s_runtime.Object) filepath.WalkFunc {
				return func(path string, info os.FileInfo, err error) error {
					return errors.New("test")
				}
			}},
			"tmp",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				templateDirectory: tt.fields.templateDirectory,
				parser:            tt.fields.parser,
				walkFunc:          tt.fields.walkFunc,
			}
			got, err := p.GetObjects(tt.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("provider.GetObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("provider.GetObjects() = %v, want %v", got, tt.want)
			}
		})
	}
}
