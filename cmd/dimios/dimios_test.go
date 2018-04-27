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
	"net/http"
	"net"
	"fmt"
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
	Context("when asked for version", func() {
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
	Context("when called with parameter webhook", func() {
		args := []string{"-webhook"}
		It("does respond with statuscode 200", func() {
			port, err := freePort()
			Expect(err).To(BeNil())
			args = append(args, fmt.Sprintf("-port=%d", port))
			serverSession, err = gexec.Start(exec.Command(pathToServerBinary, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			waitUntilPortIsOpen(port, time.Second)
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
	Context("when called without parameter webhook", func() {
		args := []string{}
		It("webserver will not be started", func() {
			port, err := freePort()
			Expect(err).To(BeNil())
			args = append(args, fmt.Sprintf("-port=%d", port))
			serverSession, err = gexec.Start(exec.Command(pathToServerBinary, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			waitUntilPortIsOpen(port, time.Second)
			_, err = http.Get(fmt.Sprintf("http://localhost:%d", port))
			Expect(err).NotTo(BeNil())
		})
	})
})

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Test Suite")
}

func waitUntilPortIsOpen(port int, maxWait time.Duration) {
	timeout := time.After(maxWait)
	for {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.IP{0, 0, 0, 0}, Port: port})
		if err != nil {
			select {
			case <-timeout:
				return
			case <-time.After(100 * time.Millisecond):
				continue
			}
		}
		conn.Close()
		return
	}
}

func freePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
