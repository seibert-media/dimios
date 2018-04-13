// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package k8s

import (
	"fmt"
	"reflect"

	k8s_appsv1 "k8s.io/api/apps/v1"
	k8s_corev1 "k8s.io/api/core/v1"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

func ObjectToString(object k8s_runtime.Object) string {
	switch t := object.(type) {
	case *k8s_appsv1.Deployment:
		return fmt.Sprintf("%s %s %s", typeof(object), t.Namespace, t.Name)
	case *k8s_corev1.Namespace:
		return fmt.Sprintf("%s %s %s", typeof(object), t.Name, t.Name)
	}
	return typeof(object)
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
