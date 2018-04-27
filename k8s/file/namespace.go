// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"os"
)

// NamespaceDirectory contains all manifest files for the namespace
type NamespaceDirectory string

// Returns the path of the NamespaceDirectory
func (n NamespaceDirectory) String() string {
	return string(n)
}

// Exists returns true if the NamespaceDirectory exists
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
