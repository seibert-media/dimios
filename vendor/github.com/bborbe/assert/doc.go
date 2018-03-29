// Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package assert provide assertions for testing

	import (
		. "github.com/bborbe/assert"
		"testing"
	)

	func TestNotNilValue(t *testing.T) {
		value := ...
		if err := AssertThat(value, NotNilValue()); err != nil {
			t.Fatal(err)
		}
	}

	func TestNilNull(t *testing.T) {
		value := ...
		if err := AssertThat(value, NilValue()); err != nil {
			t.Fatal(err)
		}
	}

	func TestEquals(t *testing.T) {
		value := ...
		if err := AssertThat(value, Is("foo")); err != nil {
			t.Fatal(err)
		}
	}

*/
package assert
