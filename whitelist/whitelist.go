// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitelist

import (
	"strings"

	k8s_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

// Entry of the whitelist.
type Entry string

func (e Entry) String() string {
	return string(e)
}

// Equals return true if name is equals.
func (e Entry) Equals(a Entry) bool {
	return e.String() == a.String()
}

// List of whitelist entries.
type List []Entry

// ByString create a whitelist from comma seperated string.
func ByString(list string) List {
	var result List
	for _, element := range strings.Split(list, ",") {
		if len(element) > 0 {
			result = append(result, Entry(element))
		}
	}
	return result
}

// IsEmpty return list is empty.
func (l List) IsEmpty() bool {
	return len(l) == 0
}

// Filter out all kubernetes object kinds not in whitelist.
func (l List) Filter(k8sobjects []k8s_runtime.Object) []k8s_runtime.Object {
	// don't filter if whitelist is empty
	if l.IsEmpty() {
		return k8sobjects
	}
	var filtered = []k8s_runtime.Object{}
	for _, object := range k8sobjects {
		u, err := k8s_runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		if err != nil {
			continue
		}
		obj := &k8s_unstructured.Unstructured{
			Object: u,
		}
		if l.Contains(Entry(obj.GetKind())) {
			filtered = append(filtered, obj)
		}
	}
	return filtered
}

// Contains return true if entry is found.
func (l List) Contains(e Entry) bool {
	for _, a := range l {
		if a.Equals(e) {
			return true
		}
	}
	return false
}
