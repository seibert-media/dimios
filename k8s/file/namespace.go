// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package file

import (
	"os"
)

// Directory contain all manifests files for the namespace.
type NamespaceDirectory string

// Returns the path of the directory.
func (n NamespaceDirectory) String() string {
	return string(n)
}

// Return true if the directory exists.
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
