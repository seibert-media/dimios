// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package file

import (
	"path"

	io_util "github.com/bborbe/io/util"
	"github.com/seibert-media/k8s-deploy/k8s"
)

// Root directory for all namespaces
type TemplateDirectory string

// Returns the path.
func (t TemplateDirectory) String() string {
	return string(t)
}

// Returns replace ~/ with the homedir.
func (d TemplateDirectory) NormalizePath() (TemplateDirectory, error) {
	root, err := io_util.NormalizePath(d.String())
	if err != nil {
		return "", err
	}
	return TemplateDirectory(root), nil
}

// Returns the NamespaceDirectory.
func (t *TemplateDirectory) PathToNamespace(namespace k8s.Namespace) NamespaceDirectory {
	return NamespaceDirectory(path.Join(t.String(), namespace.String()))
}
