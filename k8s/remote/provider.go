package remote_provider

import (
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/k8s-deploy/k8s"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	k8s_discovery "k8s.io/client-go/discovery"
	k8s_dynamic "k8s.io/client-go/dynamic"
	k8s_restclient "k8s.io/client-go/rest"
)

type provider struct {
	config *k8s_restclient.Config
}

func New(config *k8s_restclient.Config) k8s.Provider {
	return &provider{
		config: config,
	}
}

func (p *provider) GetObjects(namespace k8s.Namespace) ([]k8s_runtime.Object, error) {
	var result []k8s_runtime.Object

	glog.V(4).Infof("get objects from k8s for namespace %s", namespace)
	discoveryClient, err := k8s_discovery.NewDiscoveryClientForConfig(p.config)
	if err != nil {
		return nil, errors.Wrap(err, "creating k8s_discovery client failed")
	}
	dynamicClientPool := k8s_dynamic.NewDynamicClientPool(p.config)

	apiGroupResources, err := k8s_discovery.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, errors.Wrap(err, "get api group resouces failed")
	}
	for _, apiGroupResource := range apiGroupResources {
		for _, groupVersion := range apiGroupResource.Group.Versions {
			apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion.GroupVersion)
			if err != nil {
				return nil, errors.Wrap(err, "get server resouces failed")
			}
			client, err := dynamicClientPool.ClientForGroupVersionKind(apiResourceList.GroupVersionKind())
			if err != nil {
				return nil, errors.Wrap(err, "get client for group")
			}
			for _, resource := range apiResourceList.APIResources {
				// TODO ignore resouces with slash
				if strings.Index(resource.Name, "/") != -1 {
					continue
				}
				ri := client.Resource(&resource, namespace.String())
				object, err := ri.List(k8s_metav1.ListOptions{})
				if err != nil {
					glog.V(4).Infof("list %s failed: %s", resource.Name, err)
					continue
				}
				glog.V(4).Infof("add object %v", object)
				result = append(result, object)
			}
		}
	}
	glog.V(1).Infof("read api completed. found %d objects in namespace %s", len(result), namespace)
	return result, nil
}
