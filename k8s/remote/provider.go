// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/dimios/k8s"
	k8s_meta "k8s.io/apimachinery/pkg/api/meta"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	k8s_schema "k8s.io/apimachinery/pkg/runtime/schema"
	k8s_discovery "k8s.io/client-go/discovery"
	k8s_dynamic "k8s.io/client-go/dynamic"
)

type provider struct {
	discoveryClient   *k8s_discovery.DiscoveryClient
	dynamicClientPool k8s_dynamic.ClientPool
}

// New remote provider with passed in rest config
func New(
	discoveryClient *k8s_discovery.DiscoveryClient,
	dynamicClientPool k8s_dynamic.ClientPool,
) k8s.Provider {
	return &provider{
		discoveryClient:   discoveryClient,
		dynamicClientPool: dynamicClientPool,
	}
}

// GetObjects in the given namespace
func (p *provider) GetObjects(namespace k8s.Namespace) ([]k8s_runtime.Object, error) {
	var result []k8s_runtime.Object
	glog.V(4).Infof("get objects from k8s for namespace %s", namespace)

	resources, err := p.discoveryClient.ServerResources()
	if err != nil {
		glog.V(4).Infof("get discovery client failed: %v", err)
		return nil, errors.Wrap(err, "get server resources failed")
	}

	glog.V(4).Infof("found %d resouces in namespace %s", len(resources), namespace)

	resources = k8s_discovery.FilteredBy(
		k8s_discovery.ResourcePredicateFunc(func(groupVersion string, r *k8s_metav1.APIResource) bool {
			return k8s_discovery.SupportsAllVerbs{Verbs: []string{"list", "create"}}.Match(groupVersion, r)
		}),
		resources,
	)

	glog.V(4).Infof("keep %d resouces after filter", len(resources))

	handeled := make(map[string]struct{})

	for _, list := range resources {
		glog.V(4).Infof("list group version %v", list.GroupVersion)
		for _, resource := range list.APIResources {
			glog.V(4).Infof("check kind %v", resource.Kind)
			if _, ok := handeled[resource.Kind]; ok {
				glog.V(4).Infof("kind %v already check => skip", resource.Kind)
				continue
			}
			handeled[resource.Kind] = struct{}{}

			groupVersion, err := k8s_schema.ParseGroupVersion(list.GroupVersion)
			if err != nil {
				return nil, errors.Wrapf(err, "parse group version %s failed", list.GroupVersion)
			}
			groupVersionKind := groupVersion.WithKind(resource.Name)
			if err != nil {
				return nil, errors.Wrapf(err, "get group version for kind %s failed", resource.Name)
			}

			client, err := p.dynamicClientPool.ClientForGroupVersionKind(groupVersionKind)
			if err != nil {
				return nil, errors.Wrap(err, "get client for group")
			}

			ri := client.Resource(&resource, namespace.String())

			unstructuredList, err := ri.List(k8s_metav1.ListOptions{})
			if err != nil {
				glog.V(4).Infof("list failed: %v", err)
				continue
			}

			items, err := k8s_meta.ExtractList(unstructuredList)
			if err != nil {
				glog.V(4).Infof("extract items failed: %v", err)
				continue
			}

			glog.V(4).Infof("found %d items in kind %s", len(items), resource.Kind)

			for _, item := range items {
				glog.V(6).Infof("found api object %s", k8s.ObjectToString(item))
				is, err := IsManaged(namespace, item)
				glog.V(4).Infof("object %s ismanaged = %v", k8s.ObjectToString(item), is)
				if err != nil {
					return nil, errors.Wrap(err, "failed to determine managed state")
				}
				if is {
					glog.V(4).Infof("add %s to result", k8s.ObjectToString(item))
					result = append(result, item)
				}
			}
		}
	}
	glog.V(1).Infof("read api completed. found %d objects in namespace %s", len(result), namespace)
	return result, nil
}

// IsManaged by dimios
func IsManaged(namespace k8s.Namespace, object k8s_runtime.Object) (bool, error) {
	u, err := k8s_runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return false, errors.New("unable to convert object to k8s_unstructured")
	}
	obj := &k8s_unstructured.Unstructured{
		Object: u,
	}

	if obj.GetKind() == "Namespace" {
		if obj.GetName() != namespace.String() {
			glog.V(4).Infof("namespace missmatch")
			return false, nil
		}
	}

	if strings.HasPrefix(obj.GetKind(), "ClusterRole") {
		if obj.GetName() == "system:kube-dns-autoscaler" {
			glog.V(4).Infof("system:kube-dns-autoscaler")
			return false, nil
		}
	}

	invalidKinds := []string{"Node", "Endpoints", "CertificateSigningRequest", "Event", "ServiceAccount"}
	for _, kind := range invalidKinds {
		if obj.GetKind() == kind {
			glog.V(4).Infof("kind in %v", invalidKinds)
			return false, nil
		}
	}

	if len(obj.GetOwnerReferences()) > 0 {
		glog.V(4).Infof("object has more than 0 owner")
		return false, nil
	}

	// Labels
	if _, ok := obj.GetLabels()["kubernetes.io/bootstrapping"]; ok {
		glog.V(4).Infof("has label kubernetes.io/bootstrapping")
		return false, nil
	}
	if _, ok := obj.GetLabels()["kubernetes.io/cluster-service"]; ok {
		glog.V(4).Infof("has label kubernetes.io/cluster-service")
		return false, nil
	}
	if _, ok := obj.GetLabels()["kube-aggregator.kubernetes.io/automanaged"]; ok {
		glog.V(4).Infof("has label kube-aggregator.kubernetes.io/automanaged")
		return false, nil
	}

	// Annotations
	if _, ok := obj.GetAnnotations()["pv.kubernetes.io/provisioned-by"]; ok {
		glog.V(4).Infof("has label pv.kubernetes.io/provisioned-by")
		return false, nil
	}
	if v, ok := obj.GetAnnotations()["kubernetes.io/service-account.name"]; ok && v == "default" {
		glog.V(4).Infof("has label kubernetes.io/service-account.name")
		return false, nil
	}

	glog.V(4).Infof("is managed by dimios")
	return true, nil
}
