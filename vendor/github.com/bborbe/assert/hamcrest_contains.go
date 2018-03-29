package assert

import (
	"reflect"
	"strings"
)

type containsMatcher struct {
	expectedValue string
	message       string
}

// Contains checks if the value contains with the expected value
func Contains(expectedValue string) *containsMatcher {
	matcher := new(containsMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *containsMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *containsMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		text := value.(string)
		return strings.Contains(text, m.expectedValue)
	}
	return false
}

func (m *containsMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> contains <%v>", m.message, value, m.expectedValue)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)

}
