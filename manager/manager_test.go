// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manager

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestValidateReturnError(t *testing.T) {
	m := &Manager{}
	if err := AssertThat(m.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
