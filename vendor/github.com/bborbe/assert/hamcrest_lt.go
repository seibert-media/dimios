package assert

import (
	"reflect"
)

type ltMatcher struct {
	expectedValue interface{}
	message       string
}

// Lt checks if the value is less the expected value
func Lt(expectedValue interface{}) *ltMatcher {
	matcher := new(ltMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *ltMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *ltMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		return less(value, m.expectedValue)
	}
	return false
}

func (m *ltMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> is less than <%v>", m.message, m.expectedValue, value)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)
}
