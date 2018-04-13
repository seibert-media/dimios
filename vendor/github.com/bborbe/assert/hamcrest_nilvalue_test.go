package assert

import "testing"

func TestNilValue(t *testing.T) {
	{
		err := AssertThat(nil, NilValue())
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(t, NilValue())
		if err == nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat("", NilValue())
		if err == nil {
			t.Fatal("expect not nil")
		}
	}
	{
		err := AssertThat(0, NilValue())
		if err == nil {
			t.Fatal("expect not nil")
		}
	}
	{
		var foo *string
		err := AssertThat(foo, NilValue())
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(make([]byte, 0), NilValue())
		expectedValue := "expected nil but: was <[]>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(make([]byte, 0), NilValue().Message("msg"))
		expectedValue := "msg, expected nil but: was <[]>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}

func TestNilValuePointer(t *testing.T) {
	var value *string
	if err := AssertThat(value, NilValue()); err != nil {
		t.Fatal(err)
	}
}
