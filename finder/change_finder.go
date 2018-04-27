// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package finder

import (
	"context"
	"strings"

	"github.com/bborbe/run"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/dimios/change"
	"github.com/seibert-media/dimios/k8s"
	"github.com/seibert-media/dimios/whitelist"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Finder is looking for differences between the local file and the remote provider in `namespaces`
type Finder struct {
	FileProvider   k8s.Provider
	RemoteProvider k8s.Provider
	Namespaces     []k8s.Namespace
	Whitelist      whitelist.List
}

// New finder
func New(
	file k8s.Provider,
	remote k8s.Provider,
	namespaces []k8s.Namespace,
	whitelist whitelist.List,
) *Finder {
	return &Finder{
		FileProvider:   file,
		RemoteProvider: remote,
		Namespaces:     namespaces,
		Whitelist:      whitelist,
	}
}

// Run writes all differences found to the channel until itself or context is done
func (f *Finder) Run(ctx context.Context, c chan<- change.Change) error {
	var list []run.RunFunc
	for _, namespace := range f.Namespaces {
		n := namespace
		list = append(list, func(ctx context.Context) error {
			return f.changesForNamespace(ctx, c, n)
		})
	}
	return run.CancelOnFirstError(ctx, list...)
}

func (f *Finder) changesForNamespace(ctx context.Context, c chan<- change.Change, namespace k8s.Namespace) error {
	fileObjects, err := f.FileProvider.GetObjects(namespace)
	if err != nil {
		return errors.Wrapf(err, "get file objects failed for namespace %s", namespace)
	}

	remoteObjects, err := f.RemoteProvider.GetObjects(namespace)
	if err != nil {
		return errors.Wrapf(err, "get remote objects failed for namespace %s", namespace)
	}

	glog.V(4).Infof("found %d file objects", len(fileObjects))
	fileObjects = f.Whitelist.Filter(fileObjects)
	glog.V(4).Infof("keep %d file objects after filter", len(fileObjects))

	glog.V(4).Infof("found %d remote objects", len(remoteObjects))
	remoteObjects = f.Whitelist.Filter(remoteObjects)
	glog.V(4).Infof("keep %d remote objects after filter", len(remoteObjects))

	glog.V(4).Infof("send changes to channel")
	changeList := changes(fileObjects, remoteObjects)
	glog.V(1).Infof("got %d changes to apply for namespace %s", len(changeList), namespace)
	for _, change := range changeList {
		if writeChangeOrCancel(ctx, c, change) {
			glog.V(2).Infof("write change to channel canceled for namespace %s", namespace)
			return nil
		}
	}
	glog.V(4).Infof("all changes sent")
	return nil
}

func writeChangeOrCancel(ctx context.Context, c chan<- change.Change, change change.Change) bool {
	select {
	case c <- change:
		if glog.V(6) {
			glog.Infof("added %#v to channel", change.Object)
		} else if glog.V(4) {
			glog.Infof("added %s to channel", change.Object.GetObjectKind().GroupVersionKind().Kind)
		}
	case <-ctx.Done():
		glog.V(3).Infoln("context done, skip add changes")
		return true
	}
	return false
}

func changes(fileObjects, remoteObjects []runtime.Object) []change.Change {
	var result []change.Change
	result = append(result, deletions(fileObjects, remoteObjects)...)
	result = append(result, additions(fileObjects)...)
	return result
}

func deletions(fileObjects, remoteObjects []runtime.Object) []change.Change {
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

func additions(fileObjects []runtime.Object) []change.Change {
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
	var a1, b1 k8s_metav1.Object
	switch ta := a.(type) {
	case k8s_metav1.Object:
		a1 = ta
	}
	switch tb := b.(type) {
	case k8s_metav1.Object:
		b1 = tb
	}
	if a1 == nil || b1 == nil {
		return false
	}
	glog.V(6).Infof("name %s <> %s", a1.GetName(), b1.GetName())
	return strings.Compare(a1.GetName(), b1.GetName()) == 0
}
