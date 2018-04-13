package assert

import "testing"

func TestIsInt(t *testing.T) {
	{
		err := AssertThat(1, Is(1))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(1, Is(2))
		if err == nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat(1, Is(2))
		expectedValue := "expected <2> but got <1>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat(1, Is(2).Message("msg"))
		expectedValue := "msg, expected <2> but got <1>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}

func TestIsString(t *testing.T) {
	{
		err := AssertThat("a", Is("a"))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat("a", Is("b"))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat("a", Is("b"))
		expectedValue := "expected <b> but got <a>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
	{
		err := AssertThat("a", Is("b").Message("msg"))
		expectedValue := "msg, expected <b> but got <a>"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}

func TestIsTypeMissmatch(t *testing.T) {
	{
		err := AssertThat([]byte{}, Is(""))
		if err == nil {
			t.Fatal("expect error")
		}
	}
	{
		err := AssertThat([]byte{}, Is(""))
		expectedValue := "expected type string but got []uint8"
		if err.Error() != expectedValue {
			t.Fatalf("error message missmatch, expected '%v' but was '%v'", expectedValue, err.Error())
		}
	}
}

func TestIsByteArrayMatch(t *testing.T) {
	{
		err := AssertThat([]byte("a"), Is([]byte("a")))
		if err != nil {
			t.Fatal("expect nil")
		}
	}
	{
		err := AssertThat([]byte("a"), Is([]byte("b")))
		if err == nil {
			t.Fatal("expect error")
		}
	}
}
