package filter

import (
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	k8s_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Filter(whitelist []string, k8sobjects []k8s_runtime.Object) ([]k8s_runtime.Object) {
	var filteredk8sobjects = []k8s_runtime.Object{}

	for _, object := range k8sobjects {
		u, err := k8s_runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		if err != nil {
			continue
		}
		obj := &k8s_unstructured.Unstructured{
			Object: u,
		}
		if (contains(whitelist, obj.GetKind())) {
			filteredk8sobjects = append(filteredk8sobjects, obj)
		}
	}
	return filteredk8sobjects
}


func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}