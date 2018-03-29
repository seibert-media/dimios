package assert

import (
	"reflect"
)

type eqMatcher struct {
	expectedValue interface{}
	message       string
}

// Eq checks if the value eq with the expected value
func Eq(expectedValue interface{}) *eqMatcher {
	matcher := new(eqMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *eqMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *eqMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		return !less(m.expectedValue, value) && !less(value, m.expectedValue)
	}
	return false
}

func (m *eqMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> is equal <%v>", m.message, value, m.expectedValue)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)
}
