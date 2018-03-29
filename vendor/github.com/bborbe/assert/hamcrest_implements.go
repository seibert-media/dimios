package assert

import (
	"reflect"
)

type implementsMatcher struct {
	expectedValue interface{}
	message       string
}

// Implements checks if the value implements the expected value
func Implements(expectedValue interface{}) *implementsMatcher {
	m := new(implementsMatcher)
	m.expectedValue = expectedValue
	return m
}

func (m *implementsMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *implementsMatcher) Matches(value interface{}) bool {
	expectedType := reflect.TypeOf(m.expectedValue).Elem()
	valueType := reflect.TypeOf(value)
	return valueType.Implements(expectedType)
}

func (m *implementsMatcher) DescribeMismatch(value interface{}) error {
	expectedType := reflect.TypeOf(m.expectedValue).Elem()
	valueType := reflect.TypeOf(value)
	return buildError("expected type '%v' but got '%v'", m.message, expectedType, valueType)
}
