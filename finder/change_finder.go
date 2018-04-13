package finder

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/bborbe/k8s_deploy/change"
	"github.com/bborbe/k8s_deploy/k8s"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Finder struct {
	FileProvider   k8s.Provider
	RemoveProvider k8s.Provider
	Namespace      k8s.Namespace
}

func (f *Finder) Changes(ctx context.Context, c chan<- change.Change) error {
	defer close(c)
	fileObjects, err := f.FileProvider.GetObjects(f.Namespace)
	if err != nil {
		return fmt.Errorf("get file objects failed: %v", err)
	}
	remoteObjects, err := f.RemoveProvider.GetObjects(f.Namespace)
	if err != nil {
		return fmt.Errorf("get remote objects failed: %v", err)
	}
	for _, change := range changes(fileObjects, remoteObjects) {
		select {
		case c <- change:
			glog.V(4).Infof("added %v to channel", change)
		case <-ctx.Done():
			glog.V(3).Infoln("context done, skip add changes")
			return nil
		}
	}
	return nil
}

func changes(fileObjects, remoteObjects []runtime.Object) []change.Change {
	var result []change.Change
	result = append(result, deleteChanges(fileObjects, remoteObjects)...)
	result = append(result, applyChanges(fileObjects)...)
	glog.V(1).Infof("got %d changes to apply", len(result))
	return result
}

func deleteChanges(fileObjects, remoteObjects []runtime.Object) []change.Change {
	var result []change.Change
	for _, remoteObject := range remoteObjects {
		missing := true
		if missing && existsIn(remoteObject, fileObjects) {
			missing = false
		}
		if missing {
			result = append(result, change.Change{
				Deleted: true,
				Object:  remoteObject,
			})
		}
	}
	return result
}

func applyChanges(fileObjects []runtime.Object) []change.Change {
	var result []change.Change
	for _, fileObject := range fileObjects {
		result = append(result, change.Change{
			Object: fileObject,
		})
	}
	return result
}

func existsIn(search runtime.Object, list []runtime.Object) bool {
	for _, object := range list {
		if compare(object, search) {
			return true
		}
	}
	return false
}

func compare(a, b runtime.Object) bool {
	if a == b {
		return true
	}
	if typeof(a) != typeof(b) {
		return false
	}
	var a1, b1 metav1.ObjectMetaAccessor
	switch ta := a.(type) {
	case metav1.ObjectMetaAccessor:
		a1 = ta
	}
	switch tb := b.(type) {
	case metav1.ObjectMetaAccessor:
		b1 = tb
	}
	if a1 == nil || b1 == nil {
		return false
	}
	return strings.Compare(a1.GetObjectMeta().GetName(), b1.GetObjectMeta().GetName()) == 0
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
