package k8s

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNamespaceString(t *testing.T) {
	n := Namespace("test")
	if n.String() != "test" {
		t.Error("string conversion failed")
	}
}

func TestNamespacesFromCommaSeperatedList(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []Namespace
	}{
		{"empty", "", []Namespace{""}},
		{"one", "test", []Namespace{"test"}},
		{"two", "a,b", []Namespace{"a", "b"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := NamespacesFromCommaSeperatedList(tc.input)
			if err := AssertThat(len(list), Is(len(tc.expected))); err != nil {
				t.Fatal(err)
			}
			for i, e := range list {
				if err := AssertThat(e, Is(tc.expected[i])); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
