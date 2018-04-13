package assert

import "testing"

func TestContains(t *testing.T) {
	{
		err := AssertThat("hello world", Contains("hello"))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat("hello world", Contains("world"))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat("hello world", Contains("foo bar"))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat("hello world", Contains("foo bar"))
		expectedValue := "expected <hello world> contains <foo bar>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat("hello world", Contains("foo bar").Message("msg"))
		expectedValue := "msg, expected <hello world> contains <foo bar>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}

func TestContainsTypeMissmatch(t *testing.T) {
	{
		err := AssertThat([]byte{}, Contains("world"))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat([]byte{}, Contains("world"))
		expectedValue := "expected type string but got []uint8"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
