package assert

import "testing"

func TestLt(t *testing.T) {
	{
		err := AssertThat(3, Lt(4))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(4, Lt(3))
		if err == nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(4, Lt(3))
		expectedValue := "expected <3> is less than <4>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(4, Lt(3).Message("msg"))
		expectedValue := "msg, expected <3> is less than <4>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(4.1, Lt(3))
		expectedValue := "expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(4.1, Lt(3).Message("msg"))
		expectedValue := "msg, expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
