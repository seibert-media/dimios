package filter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Filter", func() {
	var (
		Whitelist = []string{""}
		K8sobjects = []k8s_runtime.Object{}
	)

	Describe("Whitelistfilter", func() {
		It("returns correct count of k8s objects", func() {
			var k8sdeployobjects = Filter(Whitelist, K8sobjects)
			Expect(k8sdeployobjects).ToNot(BeNil())
		})
	})
})
