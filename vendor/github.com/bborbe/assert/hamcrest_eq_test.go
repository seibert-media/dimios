package assert

import "testing"

func TestEq(t *testing.T) {
	{
		err := AssertThat(2, Eq(2))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(3, Eq(2))
		if err == nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(3, Eq(2))
		expectedValue := "expected <3> is equal <2>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(3, Eq(2).Message("msg"))
		expectedValue := "msg, expected <3> is equal <2>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(2.1, Eq(3))
		expectedValue := "expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(2.1, Eq(3).Message("msg"))
		expectedValue := "msg, expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
