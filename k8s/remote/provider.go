package remote_provider

import (
	"fmt"

	"github.com/bborbe/k8s_deploy/k8s"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
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

	ns, err := p.clientset.CoreV1().Namespaces().Get(namespace.String(), v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get namespace failed: %v", err)
	}
	result = append(result, ns)

	deploymentList, err := p.clientset.AppsV1().Deployments(namespace.String()).List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list deployments failed: %v", err)
	}
	for _, d := range deploymentList.Items {
		var obj = &d
		glog.V(2).Infof("found remote object %s", k8s.ObjectToString(obj))
		result = append(result, obj)
	}
	glog.V(1).Infof("read remote completed. found %d objects", len(result))
	return result, nil
}
