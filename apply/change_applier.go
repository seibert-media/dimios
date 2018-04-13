package apply

import (
	"context"
	"fmt"

	"github.com/bborbe/k8s_deploy/change"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
)

// Applier for changes
type Applier struct {
	config    *restclient.Config
	dynamic   dynamic.ClientPool
	discovery *discovery.DiscoveryClient
}

// New Applier with clientset
func New(config *restclient.Config) (*Applier, error) {
	discovery, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "creating discovery client failed")
	}

	return &Applier{
		config:    config,
		dynamic:   dynamic.NewDynamicClientPool(config),
		discovery: discovery,
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
			glog.V(3).Infof("apply change %v", v)
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

	client, err := c.dynamic.ClientForGroupVersionKind(change.Object.GetObjectKind().GroupVersionKind())
	if err != nil {
		return errors.Wrap(err, "creating dynamic client failed")
	}

	converter := runtime.DefaultUnstructuredConverter
	u, err := converter.ToUnstructured(change.Object)
	if err != nil {
		return errors.New("unable to convert object to unstructured")
	}
	obj := &unstructured.Unstructured{
		Object: u,
	}

	resource, err := c.getResource(client, obj)
	if err != nil {
		return errors.Wrap(err, "unable to get resource")
	}

	var result *unstructured.Unstructured
	if change.Deleted {
		err := resource.Delete(obj.GetName(), &metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrap(err, "unable to delete object")
		}
		return nil
	}
	if _, err := resource.Get(obj.GetName(), metav1.GetOptions{}); err != nil {
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

func (c *Applier) getResource(client dynamic.Interface, obj *unstructured.Unstructured) (dynamic.ResourceInterface, error) {
	res, err := c.discovery.ServerResourcesForGroupVersion(
		obj.GroupVersionKind().GroupVersion().String())
	if err != nil {
		return nil, fmt.Errorf("unable to get resources(%v) from discovery client", obj)
	}

	var resource dynamic.ResourceInterface
	for _, r := range res.APIResources {
		if r.Kind == obj.GetObjectKind().GroupVersionKind().Kind {
			resource = client.Resource(&r, obj.GetNamespace())
			return resource, nil
		}
	}

	return nil, errors.New("no ressource found for object")
}
