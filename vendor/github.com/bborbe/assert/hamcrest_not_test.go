package assert

import "testing"

func TestNotImplementsMatcher(t *testing.T) {
	m := Not(NilValue())
	var expected *Matcher
	err := AssertThat(m, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotMatches(t *testing.T) {
	err := AssertThat(t, Not(NilValue()))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotDescription(t *testing.T) {
	err := AssertThat(t, Not(NotNilValue()))
	if err == nil {
		t.Fatal("expect not nil")
	}
	err = AssertThat(err.Error(), Is("not expected not nil value"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotDescriptionWithMessage(t *testing.T) {
	err := AssertThat(t, Not(NotNilValue()).Message("foo"))
	if err == nil {
		t.Fatal("expect not nil")
	}
	err = AssertThat(err.Error(), Is("foo, not expected not nil value"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotDescriptionWithMessageMessage(t *testing.T) {
	err := AssertThat(t, Not(NotNilValue().Message("bar")).Message("foo"))
	if err == nil {
		t.Fatal("expect not nil")
	}
	err = AssertThat(err.Error(), Is("foo, not bar, expected not nil value"))
	if err != nil {
		t.Fatal(err)
	}
}
