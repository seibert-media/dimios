package errors

import "bytes"

type errorList []error

func New(errors ...error) errorList {
	return errorList(errors)
}

func NewByChan(errors <-chan error) errorList {
	var list []error
	for err := range errors {
		list = append(list, err)
	}
	return New(list...)
}

func (e errorList) Error() string {
	buf := bytes.NewBufferString("errors: ")
	first := true
	for _, err := range e {
		if first {
			first = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}
