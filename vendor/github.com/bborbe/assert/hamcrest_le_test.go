package assert

import "testing"

func TestLe(t *testing.T) {
	{
		err := AssertThat(2, Le(3))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(3, Le(3))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(3, Le(2))
		if err == nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(3, Le(2))
		expectedValue := "expected <3> is less or equal than <2>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(3, Le(2).Message("msg"))
		expectedValue := "msg, expected <3> is less or equal than <2>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(2.1, Le(3))
		expectedValue := "expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(2.1, Le(3).Message("msg"))
		expectedValue := "msg, expected type int but got float64"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}
