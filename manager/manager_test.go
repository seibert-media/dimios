// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manager_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/dimios/manager"
)

var _ = Describe("Manager", func() {
	m := &manager.Manager{}
	Context("valid initialized", func() {
		BeforeEach(func() {
			m.Kubeconfig = "/tmp/kubeconfig"
			m.Namespaces = "testns"
			m.TemplateDirectory = "/tmp/templates"
			m.TeamvaultURL = "http://teamvault.example.com"
			m.TeamvaultUser = "admin"
			m.TeamvaultPassword = "S3CR3T"
		})
		It("return no error on validation", func() {
			Expect(m.Validate()).To(BeNil())
		})
		Context("without kubeconfig", func() {
			BeforeEach(func() {
				m.Kubeconfig = ""
			})
			It("return error on validation", func() {
				Expect(m.Validate()).NotTo(BeNil())
			})
			Context("with KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT", func() {
				BeforeEach(func() {
					os.Setenv("KUBERNETES_SERVICE_HOST", "host")
					os.Setenv("KUBERNETES_SERVICE_PORT", "port")
				})
				It("return no error on validation", func() {
					Expect(m.Validate()).To(BeNil())
				})
			})
		})
	})
})

func TestSync(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manager Test Suite")
}
