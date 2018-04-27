// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package change_test

import (
	"testing"

	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/seibert-media/dimios/change"
	"github.com/seibert-media/dimios/mocks"
)

var _ = Describe("Syncer", func() {

	var (
		syncer  *change.Syncer
		applier *mocks.Applier
		getter  *mocks.Getter
		err     error
	)

	BeforeEach(func() {
		applier = &mocks.Applier{}
		getter = &mocks.Getter{}
		syncer = &change.Syncer{
			Applier: applier,
			Getter:  getter,
		}
	})

	It("read from Getter", func() {
		err = syncer.Run(context.Background())
		Expect(err).To(BeNil())
	})
})

func TestSync(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sync Test Suite")
}
