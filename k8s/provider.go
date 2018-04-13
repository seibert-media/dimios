// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package k8s

import (
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

// Provider for objects
type Provider interface {
	// Get objects for the given namespace
	GetObjects(namespace Namespace) ([]k8s_runtime.Object, error)
}
