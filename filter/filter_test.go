package filter

import (
	. "github.com/onsi/ginkgo"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Filter", func() {
	var (
		Whitelist = []string{""}
		K8sobjects = []k8s_runtime.Object{}
	)

	BeforeEach(func() {
	})

	Describe("Whitelistfilter", func() {
		It("returns correct count of k8s objects", func() {
			_ = Filter(Whitelist, K8sobjects)
		})
	})
})
