package k8s

import (
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

// Provider for objects
type Provider interface {
	// Get objects for the given namespace
	GetObjects(namespace Namespace) ([]k8s_runtime.Object, error)
}
