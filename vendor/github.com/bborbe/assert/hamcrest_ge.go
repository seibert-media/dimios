package assert

import (
	"reflect"
)

type geMatcher struct {
	expectedValue interface{}
	message       string
}

// Ge checks if the value is greater or equal the expected value
func Ge(expectedValue interface{}) *geMatcher {
	matcher := new(geMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *geMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *geMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		return !less(value, m.expectedValue)
	}
	return false
}

func (m *geMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> is greater or equal than <%v>", m.message, value, m.expectedValue)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)
}
