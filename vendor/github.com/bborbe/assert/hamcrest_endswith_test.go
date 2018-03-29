package assert

import "testing"

func TestEndswith(t *testing.T) {
	{
		err := AssertThat("hello world", Endswith("world"))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat("hello world", Endswith("hello"))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat("hello world", Endswith("hello"))
		expectedValue := "expected <hello world> ends with <hello>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat("hello world", Endswith("hello").Message("msg"))
		expectedValue := "msg, expected <hello world> ends with <hello>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}

func TestEndswithTypeMissmatch(t *testing.T) {
	{
		err := AssertThat([]byte{}, Endswith("world"))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat([]byte{}, Endswith("world"))
		expectedValue := "expected type string but got []uint8"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
