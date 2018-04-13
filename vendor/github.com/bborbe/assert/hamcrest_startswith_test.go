package assert

import "testing"

func TestStartswith(t *testing.T) {
	{
		err := AssertThat("hello world", Startswith("hello"))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat("hello world", Startswith("world"))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat("hello world", Startswith("world"))
		expectedValue := "expected <hello world> starts with <world>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat("hello world", Startswith("world").Message("msg"))
		expectedValue := "msg, expected <hello world> starts with <world>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}

func TestStartswithTypeMissmatch(t *testing.T) {
	{
		err := AssertThat([]byte{}, Startswith("world"))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat([]byte{}, Startswith("world"))
		expectedValue := "expected type string but got []uint8"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
