// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestNew(t *testing.T) {
	a := New(nil, nil)
	if err := AssertThat(a, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
