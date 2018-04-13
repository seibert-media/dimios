package apply

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/k8s-deploy/change"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	k8s_schema "k8s.io/apimachinery/pkg/runtime/schema"
	k8s_discovery "k8s.io/client-go/discovery"
	k8s_dynamic "k8s.io/client-go/dynamic"
	k8s_restclient "k8s.io/client-go/rest"
)

// Applier for changes
type Applier struct {
	staging           bool
	dynamicClientPool k8s_dynamic.ClientPool
	discoveryClient   *k8s_discovery.DiscoveryClient
}

// New Applier with clientset
func New(
	staging bool,
	config *k8s_restclient.Config,
) (*Applier, error) {
	discoveryClient, err := k8s_discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "creating k8s_discovery client failed")
	}
	dynamicClientPool := k8s_dynamic.NewDynamicClientPool(config)
	return &Applier{
		staging:           staging,
		dynamicClientPool: dynamicClientPool,
		discoveryClient:   discoveryClient,
	}, nil
}

// Apply changes being sent through the inbound channel
func (c *Applier) Apply(ctx context.Context, changes <-chan change.Change) error {
	for {
		select {
		case v, ok := <-changes:
			if !ok {
				glog.V(1).Infoln("all changes applied")
				return nil
			}
			if glog.V(6) {
				glog.Infof("added %#v to channel", v.Object)
			} else if glog.V(4) {
				glog.Infof("added %s to channel", v.Object.GetObjectKind().GroupVersionKind().Kind)
			}
			if err := c.apply(ctx, v); err != nil {
				return fmt.Errorf("apply change failed: %v", err)
			}
		case <-ctx.Done():
			glog.V(3).Infoln("context done, skip apply changes")
			return nil
		}
	}
}

func (c *Applier) apply(ctx context.Context, change change.Change) error {
	if c.staging {
		if change.Deleted {
			glog.V(0).Infof("would delete k8s object => %s", change.Object.GetObjectKind().GroupVersionKind().Kind)
			return nil
		}
		glog.V(0).Infof("would apply k8s object => %s", change.Object.GetObjectKind().GroupVersionKind().Kind)
		return nil
	}

	obj, err := createUnstructured(change)
	if err != nil {
		return errors.Wrap(err, "create unstructed failed")
	}

	resource, err := c.getResource(change.Object.GetObjectKind().GroupVersionKind(), obj)
	if err != nil {
		return errors.Wrap(err, "unable to get resource")
	}

	if change.Deleted {
		glog.V(3).Infof("delete %s", obj.GetName())
		if err := resource.Delete(obj.GetName(), &k8s_metav1.DeleteOptions{}); err != nil {
			return errors.Wrap(err, "unable to delete object")
		}
		return nil
	}
	var result *k8s_unstructured.Unstructured
	if _, err := resource.Get(obj.GetName(), k8s_metav1.GetOptions{}); err != nil {
		glog.V(3).Infoln("object not present, creating")
		result, err = resource.Create(obj)
		if err != nil {
			return errors.Wrap(err, "create object failed")
		}
	} else {
		glog.V(3).Infoln("object already present, updating")
		result, err = resource.Update(obj)
		if err != nil {
			return errors.Wrap(err, "update object failed")
		}
	}

	glog.V(6).Infof("apply result: %v", result)
	return nil
}

func createUnstructured(change change.Change) (*k8s_unstructured.Unstructured, error) {
	u, err := k8s_runtime.DefaultUnstructuredConverter.ToUnstructured(change.Object)
	if err != nil {
		return nil, errors.New("unable to convert object to k8s_unstructured")
	}
	obj := &k8s_unstructured.Unstructured{
		Object: u,
	}
	return obj, nil
}

func (c *Applier) getResource(kind k8s_schema.GroupVersionKind, obj *k8s_unstructured.Unstructured) (k8s_dynamic.ResourceInterface, error) {

	client, err := c.dynamicClientPool.ClientForGroupVersionKind(kind)
	if err != nil {
		return nil, errors.Wrap(err, "creating k8s_dynamic client failed")
	}

	res, err := c.discoveryClient.ServerResourcesForGroupVersion(obj.GroupVersionKind().GroupVersion().String())
	if err != nil {
		return nil, fmt.Errorf("unable to get resources(%v) from k8s_discovery client", obj)
	}

	for _, r := range res.APIResources {
		if r.Kind == obj.GetObjectKind().GroupVersionKind().Kind {
			return client.Resource(&r, obj.GetNamespace()), nil
		}
	}

	return nil, errors.New("no ressource found for object")
}
