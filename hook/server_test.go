// Copyright 2018 The dimios authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hook_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"github.com/seibert-media/dimios/hook"
	"net/http/httptest"
	"github.com/seibert-media/dimios/mocks"
	"errors"
	"github.com/onsi/gomega/gbytes"
	"io/ioutil"
)

var _ = Describe("Server", func() {
	var handler http.Handler
	var manager *mocks.Manager

	BeforeEach(func() {
		manager = &mocks.Manager{}
		handler = &hook.Server{
			Manager: manager,
		}
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
		Expect(gbytes.BufferWithBytes(content)).To(gbytes.Say("ok"))
	})

	It("calls run function", func() {
		Expect(manager.RunCallCount()).To(Equal(0))
		handler.ServeHTTP(httptest.NewRecorder(), &http.Request{})
		Expect(manager.RunCallCount()).To(Equal(1))
	})

	It("return status code 500 if run fails", func() {
		manager.RunReturns(errors.New("banana"))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, &http.Request{})
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusInternalServerError))
	})
	It("writes error message if run fails", func() {
		manager.RunReturns(errors.New("banana"))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, &http.Request{})
		content, _ := ioutil.ReadAll(recorder.Result().Body)
		Expect(gbytes.BufferWithBytes(content)).To(gbytes.Say("run failed: banana"))
	})
})

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Test Suite")
}
