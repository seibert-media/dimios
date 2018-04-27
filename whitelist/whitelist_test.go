// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package whitelist_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/dimios/whitelist"
	"k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
	k8s_metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_runtime "k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Filter", func() {
	var (
		whitelist      = whitelist.List([]whitelist.Entry{""})
		testDeployment = v1.Deployment{
			TypeMeta: k8s_metav1.TypeMeta{
				APIVersion: "extensions/v1beta1",
				Kind:       "Deployment",
			},
		}
		testIngress = v1beta1.Ingress{
			TypeMeta: k8s_metav1.TypeMeta{
				APIVersion: "extensions/v1beta1",
				Kind:       "Ingress",
			},
		}
		k8sObjects = []k8s_runtime.Object{&testDeployment, &testIngress}
	)

	Describe("Whitelistfilter", func() {
		It("returns correct count of k8s objects with empty whitelist", func() {
			var k8sdeployobjects = whitelist.Filter(k8sObjects)
			Expect(k8sdeployobjects).ToNot(BeNil())
			Expect(len(k8sdeployobjects)).To(Equal(len(k8sdeployobjects)))
		})
		It("returns correct count of k8s objects with whitelist deployment", func() {
			whitelist = append(whitelist, "Deployment")
			var k8sdeployobjects = whitelist.Filter(k8sObjects)
			Expect(len(k8sdeployobjects)).To(Equal(1))
		})
		It("returns correct count of k8s objects with whitelist deployment and ingress", func() {
			whitelist = append(whitelist, "Deployment", "Ingress")
			var k8sdeployobjects = whitelist.Filter(k8sObjects)
			Expect(len(k8sdeployobjects)).To(Equal(2))
		})
	})
})

func TestWhitelist(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Whitelist Test Suite")
}
