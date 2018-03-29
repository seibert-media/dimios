package assert

import (
	"reflect"
)

type notNilValueMatcher struct {
	message string
}

// NotNilValue checks if the value is not nil
func NotNilValue() *notNilValueMatcher {
	matcher := new(notNilValueMatcher)
	return matcher
}

func (m *notNilValueMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *notNilValueMatcher) Matches(value interface{}) bool {
	if value == nil {
		return false
	}
	r := reflect.ValueOf(value)
	if r.Kind() != reflect.Ptr {
		return true
	}
	return !r.IsNil()
}

func (m *notNilValueMatcher) DescribeMismatch(value interface{}) error {
	return buildError("expected not nil value", m.message)
}
