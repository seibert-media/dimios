// Copyright 2018 The dimios authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hook_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/seibert-media/dimios/hook"
	"github.com/seibert-media/dimios/mocks"
)

var _ = Describe("Server", func() {
	var handler http.Handler
	var manager *mocks.Manager

	BeforeEach(func() {
		manager = &mocks.Manager{}
		handler = hook.NewHandler(manager)
	})

	It("return status code 200", func() {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, &http.Request{})
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})
	It("write ok on success", func() {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, &http.Request{})
		content, _ := ioutil.ReadAll(recorder.Result().Body)
		Expect(gbytes.BufferWithBytes(content)).To(gbytes.Say("sync triggerd"))
	})
	It("calls run function", func() {
		Expect(manager.RunCallCount()).To(Equal(0))
		handler.ServeHTTP(httptest.NewRecorder(), &http.Request{})
		time.Sleep(10 * time.Millisecond)
		Expect(manager.RunCallCount()).To(Equal(1))
	})
})

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Test Suite")
}
