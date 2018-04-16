package k8s

import (
	"testing"

	k8s_appsv1 "k8s.io/api/apps/v1"
	k8s_corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

func TestObjectToString(t *testing.T) {
	if res := ObjectToString(createIngress("test", "test")); res != "*v1beta1.Ingress" {
		t.Error("ObjectToString(ingress) unexpected result", res)
	}
	if res := ObjectToString(createDeployment("test", "test")); res != "*v1.Deployment test test" {
		t.Error("ObjectToString(deployment) unexpected result", res)
	}
	if res := ObjectToString(createNamespace("test")); res != "*v1.Namespace test" {
		t.Error("ObjectToString(namespace) unexpected result", res)
	}
}

func createIngress(ns, name string) k8s_runtime.Object {
	return &v1beta1.Ingress{
		TypeMeta: k8s_metav1.TypeMeta{
			Kind: "Ingress",
		},
		ObjectMeta: k8s_metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func createDeployment(ns, name string) k8s_runtime.Object {
	return &k8s_appsv1.Deployment{
		TypeMeta: k8s_metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: k8s_metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func createNamespace(name string) k8s_runtime.Object {
	return &k8s_corev1.Namespace{
		TypeMeta: k8s_metav1.TypeMeta{
			Kind: "Namespace",
		},
		ObjectMeta: k8s_metav1.ObjectMeta{
			Name: name,
		},
	}
}
