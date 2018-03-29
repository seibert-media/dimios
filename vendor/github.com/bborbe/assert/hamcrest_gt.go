package assert

import (
	"reflect"
)

type gtMatcher struct {
	expectedValue interface{}
	message       string
}

// Gt checks if the value is greater the expected value
func Gt(expectedValue interface{}) *gtMatcher {
	matcher := new(gtMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *gtMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *gtMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		return less(m.expectedValue, value)
	}
	return false
}

func (m *gtMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> is greater than <%v>", m.message, value, m.expectedValue)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)
}
