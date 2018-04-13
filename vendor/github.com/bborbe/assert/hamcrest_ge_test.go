package assert

import "testing"

func TestGe(t *testing.T) {
	{
		err := AssertThat(3, Ge(2))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(3, Ge(3))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(2, Ge(3))
		if err == nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(2, Ge(3))
		expectedValue := "expected <2> is greater or equal than <3>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(2, Ge(3).Message("msg"))
		expectedValue := "msg, expected <2> is greater or equal than <3>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(2.1, Ge(3))
		expectedValue := "expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(2.1, Ge(3).Message("msg"))
		expectedValue := "msg, expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
