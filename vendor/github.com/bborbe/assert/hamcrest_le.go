package assert

import (
	"reflect"
)

type leMatcher struct {
	expectedValue interface{}
	message       string
}

// Le checks if the value is less or equal the expected value
func Le(expectedValue interface{}) *leMatcher {
	matcher := new(leMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *leMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *leMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		return !less(m.expectedValue, value)
	}
	return false
}

func (m *leMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> is less or equal than <%v>", m.message, value, m.expectedValue)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)
}
