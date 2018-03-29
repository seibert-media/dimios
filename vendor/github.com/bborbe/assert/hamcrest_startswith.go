package assert

import (
	"reflect"
	"strings"
)

type startswithMatcher struct {
	expectedValue string
	message       string
}

// Startswith checks if the value starts with the expected value
func Startswith(expectedValue string) *startswithMatcher {
	matcher := new(startswithMatcher)
	matcher.expectedValue = expectedValue
	return matcher
}

func (m *startswithMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *startswithMatcher) Matches(value interface{}) bool {
	if sameType(value, m.expectedValue) {
		text := value.(string)
		return strings.Index(text, m.expectedValue) == 0
	}
	return false
}

func (m *startswithMatcher) DescribeMismatch(value interface{}) error {
	if sameType(value, m.expectedValue) {
		return buildError("expected <%v> starts with <%v>", m.message, value, m.expectedValue)
	}
	expectedType := reflect.TypeOf(m.expectedValue)
	valueType := reflect.TypeOf(value)
	return buildError("expected type %v but got %v", m.message, expectedType, valueType)

}
