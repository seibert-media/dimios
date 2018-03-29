package assert

// Matcher describe a check
type Matcher interface {
	Matches(value interface{}) bool
	DescribeMismatch(value interface{}) error
	Message(message string) Matcher
}

// AssertThat is used to compare a value with a matcher
func AssertThat(value interface{}, matcher Matcher) error {
	return That(value, matcher)
}

// That is the same as AssertThat
func That(value interface{}, matcher Matcher) error {
	if matcher.Matches(value) {
		return nil
	}
	return matcher.DescribeMismatch(value)
}
