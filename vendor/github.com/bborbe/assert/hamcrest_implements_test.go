package assert

import "testing"

type TestInterface interface {
	Print()
}

type TestObjectWithoutPrint struct{}

type TestObjectWithPrint struct{}

func (*TestObjectWithPrint) Print() {}

func TestImplements(t *testing.T) {
	{
		var i *TestInterface = nil
		value := new(TestObjectWithPrint)
		err := AssertThat(value, Implements(i))
		if err != nil {
			t.Fatal("shouldn't return error")
		}
	}
	{
		var i *TestInterface = nil
		value := new(TestObjectWithoutPrint)
		err := AssertThat(value, Implements(i))
		if err == nil {
			t.Fatal("should return error")
		}
	}
	{
		var i *TestInterface = nil
		value := new(TestObjectWithoutPrint)
		err := AssertThat(value, Implements(i))
		if err.Error() != "expected type 'assert.TestInterface' but got '*assert.TestObjectWithoutPrint'" {
			t.Fatal("errormessage is incorrect")
		}
	}
	{
		var i *TestInterface = nil
		value := new(TestObjectWithoutPrint)
		err := AssertThat(value, Implements(i).Message("msg"))
		if err.Error() != "msg, expected type 'assert.TestInterface' but got '*assert.TestObjectWithoutPrint'" {
			t.Fatal("errormessage is incorrect")
		}
	}
}
