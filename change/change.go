// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package change

import (
	"fmt"

	"github.com/seibert-media/k8s-deploy/k8s"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

// Contains the Kubernetes object and if it was deleted or not
type Change struct {
	Deleted bool
	Object  k8s_runtime.Object
}

// String representation of the change.
func (c *Change) String() string {
	if c.Deleted {
		return fmt.Sprintf("DELETE %s", k8s.ObjectToString(c.Object))
	} else {
		return fmt.Sprintf("CREATE %s", k8s.ObjectToString(c.Object))
	}
}
