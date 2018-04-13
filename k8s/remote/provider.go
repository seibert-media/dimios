package remote_provider

import (
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
	discoveryClient, err := k8s_discovery.NewDiscoveryClientForConfig(p.config)
	if err != nil {
		return nil, errors.Wrap(err, "creating k8s_discovery client failed")
	}
	dynamicClientPool := k8s_dynamic.NewDynamicClientPool(p.config)
	resourceLists, err := discoveryClient.ServerResources()
	if err != nil {
		return nil, errors.Wrap(err, "get server resouces failed")
	}
	var result []k8s_runtime.Object
	for _, resourceList := range resourceLists {
		client, err := dynamicClientPool.ClientForGroupVersionKind(resourceList.GroupVersionKind())
		if err != nil {
			return nil, errors.Wrap(err, "get client for group")
		}
		for _, resource := range resourceList.APIResources {
			ri := client.Resource(&resource, namespace.String())
			object, err := ri.List(k8s_metav1.ListOptions{})
			if err != nil {
				return nil, err
			}
			result = append(result, object)
		}
	}
	return result, nil
}
