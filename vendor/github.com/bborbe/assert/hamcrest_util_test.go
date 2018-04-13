package assert

import (
	"testing"
)

func TestBuildError(t *testing.T) {
	{
		err := buildError("expected type '%s' but got '%s'", "", "foo", "bar")
		if err == nil {
			t.Fatal("err is nil")
		}
		if err.Error() != "expected type 'foo' but got 'bar'" {
			t.Fatal("errormessage is incorrect")
		}
	}
	{
		err := buildError("expected type '%s' but got '%s'", "message", "foo", "bar")
		if err == nil {
			t.Fatal("err is nil")
		}
		if err.Error() != "message, expected type 'foo' but got 'bar'" {
			t.Fatal("errormessage is incorrect")
		}
	}
}

func TestIsByteArray(t *testing.T) {
	{
		err := AssertThat(isByteArray(nil), Is(false))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := AssertThat(isByteArray([]string{"a", "b"}), Is(false))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := AssertThat(isByteArray([]byte("hello")), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLessTypeNotEqualAlwaysReturnFalse(t *testing.T) {
	if less(int32(1), int64(2)) {
		t.Fatal("invalid")
	}
}

func TestLessNotSupportTypePanics(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		} else {
			t.Fatal("panic should be called")
		}
	}()
	less("1", "1")
}

func TestLess(t *testing.T) {
	if less(int(1), int(1)) {
		t.Fatal("invalid")
	}
	if less(int8(1), int8(1)) {
		t.Fatal("invalid")
	}
	if less(int16(1), int16(1)) {
		t.Fatal("invalid")
	}
	if less(int32(1), int32(1)) {
		t.Fatal("invalid")
	}
	if less(int64(1), int64(1)) {
		t.Fatal("invalid")
	}
	if less(uint(1), uint(1)) {
		t.Fatal("invalid")
	}
	if less(uint8(1), uint8(1)) {
		t.Fatal("invalid")
	}
	if less(uint16(1), uint16(1)) {
		t.Fatal("invalid")
	}
	if less(uint32(1), uint32(1)) {
		t.Fatal("invalid")
	}
	if less(uint64(1), uint64(1)) {
		t.Fatal("invalid")
	}
	if less(float32(1), float32(1)) {
		t.Fatal("invalid")
	}
	if less(float64(1), float64(1)) {
		t.Fatal("invalid")
	}
}
