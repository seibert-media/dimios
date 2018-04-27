package k8s

import (
	"strings"
)

// Namespace in Kubernetes
type Namespace string

// Return namespace as string.
func (n Namespace) String() string {
	return string(n)
}

// NamespacesFromList converts the string list to a namespace list.
func NamespacesFromList(namespaces []string) []Namespace {
	result := make([]Namespace, len(namespaces))
	for i, v := range namespaces {
		result[i] = Namespace(v)
	}
	return result
}

// NamespacesFromCommaSeperatedList returns a list of namespaces parsed from string
func NamespacesFromCommaSeperatedList(list string) []Namespace {
	return NamespacesFromList(strings.Split(list, ","))
}

func WhitelistFromCommaSeperatedList(list string) []string {
	return strings.Split(list, ",")
}
