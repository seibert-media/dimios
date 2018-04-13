package finder

import (
	"context"
	"reflect"
	"strings"

	"github.com/bborbe/run"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/k8s-deploy/change"
	"github.com/seibert-media/k8s-deploy/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Finder struct {
	FileProvider   k8s.Provider
	RemoveProvider k8s.Provider
	Namespaces     []k8s.Namespace
}

type changeNamespace func(context.Context) error

func (f *Finder) Changes(ctx context.Context, c chan<- change.Change) error {
	defer close(c)
	var list []run.RunFunc
	for _, namespace := range f.Namespaces {
		list = append(list, func(ctx context.Context) error {
			return f.changesForNamespace(ctx, c, namespace)
		})
	}
	return run.CancelOnFirstError(ctx, list...)
}

func (f *Finder) changesForNamespace(ctx context.Context, c chan<- change.Change, namespace k8s.Namespace) error {
	fileObjects, err := f.FileProvider.GetObjects(namespace)
	if err != nil {
		return errors.Wrap(err, "get file objects failed")
	}
	remoteObjects, err := f.RemoveProvider.GetObjects(namespace)
	if err != nil {
		return errors.Wrap(err, "get remote objects failed")
	}
	for _, change := range changes(fileObjects, remoteObjects) {
		select {
		case c <- change:
			if glog.V(6) {
				glog.Infof("added %#v to channel", change.Object)
			} else if glog.V(4) {
				glog.Infof("added %s to channel", change.Object.GetObjectKind().GroupVersionKind().Kind)
			}
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
			glog.V(2).Infof("delete %s", remoteObject.GetObjectKind().GroupVersionKind().Kind)
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
		glog.V(2).Infof("apply %s", fileObject.GetObjectKind().GroupVersionKind().Kind)
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
	glog.V(6).Infof("type %s <> %s", a.GetObjectKind().GroupVersionKind().Kind, b.GetObjectKind().GroupVersionKind().Kind)
	if a.GetObjectKind().GroupVersionKind().Kind != b.GetObjectKind().GroupVersionKind().Kind {
		return false
	}
	var a1, b1 metav1.Object
	switch ta := a.(type) {
	case metav1.Object:
		a1 = ta
	}
	switch tb := b.(type) {
	case metav1.Object:
		b1 = tb
	}
	if a1 == nil || b1 == nil {
		return false
	}
	glog.V(6).Infof("name %s <> %s", a1.GetName(), b1.GetName())
	return strings.Compare(a1.GetName(), b1.GetName()) == 0
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
