// Copyright 2018 The dimios authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const connectionTimeout = 5 * time.Second

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

type args map[string]string

func (a args) list() []string {
	var result []string
	for k, v := range a {
		if len(v) == 0 {
			result = append(result, fmt.Sprintf("-%s", k))
		} else {
			result = append(result, fmt.Sprintf("-%s=%s", k, v))
		}
	}
	return result
}

var _ = Describe("the dimios", func() {
	var err error
	validargs := args{"staging": "true"}
	Context("when asked for version", func() {
		BeforeEach(func() {
			validargs["version"] = ""
		})
		It("prints version string", func() {
			serverSession, err = gexec.Start(exec.Command(pathToServerBinary, validargs.list()...), GinkgoWriter, GinkgoWriter)
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
	Context("with valid args", func() {
		var _ = BeforeEach(func() {
			delete(validargs, "version")
			validargs["logtostderr"] = ""
			validargs["v"] = "2"
			validargs["dir"] = os.TempDir()
			validargs["namespaces"] = "testns"
			validargs["teamvault-url"] = "http://teamvault.example.com"
			validargs["teamvault-user"] = "admin"
			validargs["teamvault-pass"] = "S3CR3T"
			validargs["kubeconfig"] = "~/.kube/config"
		})
		Context("with port", func() {
			var port int
			BeforeEach(func() {
				port, err = freePort()
				Expect(err).To(BeNil())
				validargs["port"] = strconv.Itoa(port)
			})
			Context("when called without parameter webhook", func() {
				BeforeEach(func() {
					delete(validargs, "webhook")
				})
				It("webserver will not be started", func() {
					serverSession, err = gexec.Start(exec.Command(pathToServerBinary, validargs.list()...), GinkgoWriter, GinkgoWriter)
					Expect(err).To(BeNil())
					waitUntilPortIsOpen(port, 500*time.Millisecond)
					_, err = http.Get(fmt.Sprintf("http://localhost:%d", port))
					Expect(err).NotTo(BeNil())
				})
			})
			Context("when called with parameter webhook", func() {
				BeforeEach(func() {
					validargs["webhook"] = "true"
				})
				It("does respond with statuscode 200", func() {
					fmt.Printf("%v", validargs)
					serverSession, err = gexec.Start(exec.Command(pathToServerBinary, validargs.list()...), GinkgoWriter, GinkgoWriter)
					Expect(err).To(BeNil())
					waitUntilPortIsOpen(port, connectionTimeout)
					resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(BeNumerically(">", 0))
				})
			})
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
