package remote_provider

import (
	"github.com/seibert-media/k8s-deploy/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	k8s_dynamic "k8s.io/client-go/dynamic"
	k8s_discovery "k8s.io/client-go/discovery"
	"github.com/pkg/errors"
)

type provider struct {
	clientset *kubernetes.Clientset
}

func New(clientset *kubernetes.Clientset) k8s.Provider {
	return &provider{
		clientset: clientset,
	}
}

func (p *provider) GetObjects(namespace k8s.Namespace) ([]runtime.Object, error) {
	var result []runtime.Object

	var discoveryClient *k8s_discovery.DiscoveryClient
	var client k8s_dynamic.Interface

	resourceLists, err := discoveryClient.ServerResources()
	if err != nil {
		return nil, errors.Wrap(err, "get server resouces failed")
	}
	for _, resourceList := range resourceLists {
		for _, resource := range resourceList.APIResources {
			ri := client.Resource(&resource, namespace.String())
			object, err := ri.List(metav1.ListOptions{})
			if err != nil {
				return nil, err
			}
			result = append(result, object)
		}
	}
	return result, nil
}
