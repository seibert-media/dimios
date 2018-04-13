package finder

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/seibert-media/k8s-deploy/change"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestNew(t *testing.T) {
	a := &Finder{}
	if err := AssertThat(a, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func createIngress(ns, name string) runtime.Object {
	return &v1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind: "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func createDeployment(ns, name string) runtime.Object {
	return &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}
func TestCompare(t *testing.T) {
	i1 := createIngress("debug", "hello")
	i2 := createIngress("debug", "world")
	i3 := createIngress("debug", "hello")
	d1 := createDeployment("debug", "hello")
	testCases := []struct {
		name     string
		a        runtime.Object
		b        runtime.Object
		expected bool
	}{
		{"ref equal", i1, i1, true},
		{"ref not equal", i1, i2, false},
		{"ref not equal but content", i1, i3, true},
		{"name equal by kind not equal", i1, d1, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compare(tc.a, tc.b)
			if err := AssertThat(result, Is(tc.expected)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestExistsIn(t *testing.T) {
	a := createIngress("debug", "hello")
	b := createIngress("debug", "world")
	testCases := []struct {
		name     string
		search   runtime.Object
		list     []runtime.Object
		expected bool
	}{
		{"empty list", a, []runtime.Object{}, false},
		{"object in list 1", a, []runtime.Object{a}, true},
		{"object in list 2", a, []runtime.Object{a, b}, true},
		{"object in list 3", a, []runtime.Object{b, a}, true},
		{"object not in list", a, []runtime.Object{b}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := existsIn(tc.search, tc.list)
			if err := AssertThat(result, Is(tc.expected)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestApplyChanges(t *testing.T) {
	a := createIngress("debug", "hello")
	b := createIngress("debug", "world")
	testCases := []struct {
		name        string
		fileObjects []runtime.Object
		expected    []change.Change
	}{
		{"no new", []runtime.Object{}, []change.Change{}},
		{"one new", []runtime.Object{a}, []change.Change{{Deleted: false, Object: a}}},
		{"two new", []runtime.Object{a, b}, []change.Change{{Deleted: false, Object: a}, {Deleted: false, Object: b}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := applyChanges(tc.fileObjects)
			if err := AssertThat(len(result), Is(len(tc.expected)).Message("length of result mismatch")); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < len(tc.expected); i++ {
				if err := AssertThat(result[i], Is(tc.expected[i])); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestDeleteChanges(t *testing.T) {
	a := createIngress("debug", "hello")
	b := createIngress("debug", "world")
	c := createDeployment("debug", "hello")
	testCases := []struct {
		name          string
		fileObjects   []runtime.Object
		remoteObjects []runtime.Object
		expected      []change.Change
	}{
		{"empty", []runtime.Object{}, []runtime.Object{}, []change.Change{}},
		{"nothing to delete", []runtime.Object{a}, []runtime.Object{}, []change.Change{}},
		{"one to delete", []runtime.Object{}, []runtime.Object{a}, []change.Change{{Deleted: true, Object: a}}},
		{"two object, but nothing to delete", []runtime.Object{b}, []runtime.Object{a, b}, []change.Change{{Deleted: true, Object: a}}},
		{"two already exists", []runtime.Object{a, b}, []runtime.Object{b, a}, []change.Change{}},
		{"same name but different type", []runtime.Object{a}, []runtime.Object{c}, []change.Change{{Deleted: true, Object: c}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := deleteChanges(tc.fileObjects, tc.remoteObjects)
			if err := AssertThat(len(result), Is(len(tc.expected)).Message("length of result mismatch")); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < len(tc.expected); i++ {
				if err := AssertThat(result[i], Is(tc.expected[i])); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
