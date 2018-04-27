package filter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Filter", func() {
	var (
		Whitelist      = []string{""}
		TestDeployment = v1.Deployment{
			TypeMeta: k8s_metav1.TypeMeta{
				APIVersion: "extensions/v1beta1",
				Kind:       "Deployment",
			},
		}
		TestIngress = v1beta1.Ingress{
			TypeMeta: k8s_metav1.TypeMeta{
				APIVersion: "extensions/v1beta1",
				Kind:       "Ingress",
			},
		}
		K8sobjects = []k8s_runtime.Object{&TestDeployment, &TestIngress}
	)

	Describe("Whitelistfilter", func() {
		It("returns correct count of k8s objects with empty whitelist", func() {
			var k8sdeployobjects = Filter(Whitelist, K8sobjects)
			Expect(k8sdeployobjects).ToNot(BeNil())
			Expect(len(k8sdeployobjects)).To(Equal(0))
		})
		It("returns correct count of k8s objects with whitelist deployment", func() {
			Whitelist = append(Whitelist, "Deployment")
			var k8sdeployobjects = Filter(Whitelist, K8sobjects)
			Expect(len(k8sdeployobjects)).To(Equal(1))
		})
		It("returns correct count of k8s objects with whitelist deployment and ingress", func() {
			Whitelist = append(Whitelist, "Deployment", "Ingress")
			var k8sdeployobjects = Filter(Whitelist, K8sobjects)
			Expect(len(k8sdeployobjects)).To(Equal(2))
		})
	})
})
