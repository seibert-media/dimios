package k8s

import (
	"fmt"
	"reflect"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func ObjectToString(object runtime.Object) string {
	switch t := object.(type) {
	case *v1.Deployment:
		return fmt.Sprintf("%s %s %s", typeof(object), t.Namespace, t.Name)
	case *corev1.Namespace:
		return fmt.Sprintf("%s %s %s", typeof(object), t.Name, t.Name)
	}
	return typeof(object)
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
