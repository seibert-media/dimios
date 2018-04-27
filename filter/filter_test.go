package filter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/apps/v1"
)

var _ = Describe("Filter", func() {
	var (
		Whitelist = []string{""}
		TestDeployment = v1.Deployment{
			TypeMeta: k8s_metav1.TypeMeta{
				APIVersion: "extensions/v1beta1",
				Kind:       "Deployment",
				},
		}
		K8sobjects = []k8s_runtime.Object{&TestDeployment}
	)

	Describe("Whitelistfilter", func() {
		It("returns correct count of k8s objects with empty whitelist", func() {
			var k8sdeployobjects, _ = Filter(Whitelist, K8sobjects)
			Expect(k8sdeployobjects).ToNot(BeNil())
			Expect(len(k8sdeployobjects)).To(Equal(0))
		})
		It("returns correct count of k8s objects with whitelist deployment", func() {
			Whitelist = append(Whitelist, "Deployment")
			var k8sdeployobjects, _ = Filter(Whitelist, K8sobjects)
			Expect(len(k8sdeployobjects)).To(Equal(1))
		})
	})
})
