package header

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestCreateParseBearer(t *testing.T) {
	header := CreateAuthorizationBearerHeader("foo", "bar")
	name, value, err := ParseAuthorizationHeader("Bearer", header)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(name, Is("foo")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(value, Is("bar")); err != nil {
		t.Fatal(err)
	}
}
