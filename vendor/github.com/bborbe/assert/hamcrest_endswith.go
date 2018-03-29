package assert

import (
	"reflect"
	"strings"
)

type endswithMatcher struct {
	expectedValue string
	message       string
}

// Endswith checks if the value ends with the expected value
func Endswith(expectedValue string) *endswithMatcher {
	matcher := new(endswithMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *endswithMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *endswithMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		text := value.(string)
		return strings.LastIndex(text, m.expectedValue) == len(text)-len(m.expectedValue)
	}
	return false
}

func (m *endswithMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> ends with <%v>", m.message, value, m.expectedValue)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)
}
