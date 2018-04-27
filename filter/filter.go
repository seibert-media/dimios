package filter

import k8s_runtime "k8s.io/apimachinery/pkg/runtime"

func Filter(whitelist []string, k8sobjects []k8s_runtime.Object) []k8s_runtime.Object {
	return k8sobjects
}
