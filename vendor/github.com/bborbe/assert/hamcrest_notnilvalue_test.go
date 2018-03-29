package assert

import "testing"

func TestNotNilValue(t *testing.T) {
	{
		err := AssertThat(t, NotNilValue())
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(0, NotNilValue())
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat("", NotNilValue())
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(nil, NotNilValue())
		if err == nil {
			t.Fatal("expect nil")
		}
	}
	{
		var foo *string
		err := AssertThat(foo, NotNilValue())
		if err == nil {
			t.Fatal("expect not nil")
		}
	}
	{
		err := AssertThat(nil, NotNilValue())
		expectedValue := "expected not nil value"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(nil, NotNilValue().Message("msg"))
		expectedValue := "msg, expected not nil value"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
