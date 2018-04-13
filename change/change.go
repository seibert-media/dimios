// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package change

import (
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

// Change stores the Kubernetes object and whether to delete it or not
type Change struct {
	Deleted bool
	Object  k8s_runtime.Object
}
