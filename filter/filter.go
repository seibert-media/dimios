package filter

import (
	k8s_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

// Filter out all kubernetes object kinds not in whitelist
func Filter(whitelist []string, k8sobjects []k8s_runtime.Object) []k8s_runtime.Object {
	var filtered = []k8s_runtime.Object{}

	for _, object := range k8sobjects {
		u, err := k8s_runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		if err != nil {
			continue
		}
		obj := &k8s_unstructured.Unstructured{
			Object: u,
		}
		if contains(whitelist, obj.GetKind()) {
			filtered = append(filtered, obj)
		}
	}
	return filtered
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
