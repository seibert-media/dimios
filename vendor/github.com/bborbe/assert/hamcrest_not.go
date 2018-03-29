package assert

type notValueMatcher struct {
	message       string
	expectedValue Matcher
}

// Not checks if the value matches not expected matcher
func Not(expectedValue Matcher) *notValueMatcher {
	m := new(notValueMatcher)
	m.expectedValue = expectedValue
	return m
}

func (m *notValueMatcher) Message(message string) Matcher {
	m.message = message
	return m
}

func (m *notValueMatcher) Matches(value interface{}) bool {
	return !m.expectedValue.Matches(value)
}

func (m *notValueMatcher) DescribeMismatch(value interface{}) error {
	err := m.expectedValue.DescribeMismatch(value)
	return buildError("not %s", m.message, err.Error())
}
