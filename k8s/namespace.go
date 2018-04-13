package k8s

import (
	"strings"
)

// Namespace of K8s
type Namespace string

// Return namespace as string.
func (n Namespace) String() string {
	return string(n)
}

// Converts the string list into a namespace list.
func NamespacesFromList(namespaces []string) []Namespace {
	result := make([]Namespace, len(namespaces))
	for i, v := range namespaces {
		result[i] = Namespace(v)
	}
	return result
}

// Returns a list of namespaces parsed from the input.
func NamespacesFromCommaSeperatedList(list string) []Namespace {
	return NamespacesFromList(strings.Split(list, ","))
}
