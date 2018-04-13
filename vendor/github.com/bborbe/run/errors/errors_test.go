package errors

import (
	"testing"
	"fmt"
)

func TestNewEmptyError(t *testing.T) {
	err := New()
	if err != nil {
		t.Fatalf("nil expected")
	}
}

func TestNewErorrList(t *testing.T) {
	err := New(fmt.Errorf("test"))
	if err == nil {
		t.Fatalf("nil not expected")
	}
	if err.Error() != "errors: test" {
		t.Fatalf("invalid msg")
	}
}

func TestNewByChanEmptyError(t *testing.T) {
	c := make(chan error, 10)
	close(c)
	err := NewByChan(c)
	if err != nil {
		t.Fatalf("nil expected")
	}
}

func TestNewByChanErorrList(t *testing.T) {
	c := make(chan error, 10)
	c <- fmt.Errorf("test")
	close(c)
	err := NewByChan(c)
	if err == nil {
		t.Fatalf("nil not expected")
	}
	if err.Error() != "errors: test" {
		t.Fatalf("invalid msg")
	}
}
