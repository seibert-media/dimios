// Copyright 2018 The dimios authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var pathToServerBinary string
var serverSession *gexec.Session

var _ = BeforeSuite(func() {
	var err error
	pathToServerBinary, err = gexec.Build("github.com/seibert-media/dimios/cmd/dimios")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	serverSession.Interrupt()
	Eventually(serverSession).Should(gexec.Exit())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("the dimios", func() {
	var err error
	Describe("when asked for version", func() {
		It("prints version string", func() {
			serverSession, err = gexec.Start(exec.Command(pathToServerBinary, "-version"), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			serverSession.Wait(time.Second)
			Expect(serverSession.ExitCode()).To(Equal(0))
			Expect(serverSession.Out).To(gbytes.Say(`-- Dimios --
unknown - version unknown
  branch: 	unknown
  revision: 	unknown
  build date: 	unknown
  build user: 	unknown
  go version: 	unknown
`))
		})
	})
})

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Test Suite")
}
