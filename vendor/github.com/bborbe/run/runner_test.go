package run

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
	"time"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestCancelOnFirstFinishRunNothing(t *testing.T) {
	err := CancelOnFirstFinish(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCancelOnFirstFinishReturnOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.NewTicker(100 * time.Millisecond).C
		cancel()
	}()
	err := CancelOnFirstFinish(ctx,
		func(ctx context.Context) error {
			<-time.NewTicker(time.Minute).C
			return nil
		})
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCancelOnFirstFinishRun(t *testing.T) {
	r1 := new(testRunnable)
	err := CancelOnFirstFinish(context.Background(), r1.Run)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r1.counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestCancelOnFirstFinishRunThree(t *testing.T) {
	r1 := new(testRunnable)
	err := CancelOnFirstFinish(context.Background(), r1.Run, r1.Run, r1.Run)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r1.counter, Ge(1)); err != nil {
		t.Fatal(err)
	}
}

func TestCancelOnFirstFinishRunFail(t *testing.T) {
	r1 := new(testRunnable)
	r1.result = errors.New("fail")
	err := CancelOnFirstFinish(context.Background(), r1.Run)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r1.counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestCancelOnFirstErrorRunNothing(t *testing.T) {
	err := CancelOnFirstError(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestCancelOnFirstErrorReturnOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.NewTicker(100 * time.Millisecond).C
		cancel()
	}()
	err := CancelOnFirstError(ctx,
		func(ctx context.Context) error {
			<-time.NewTicker(time.Minute).C
			return nil
		})
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

type testRunnable struct {
	counter int
	result  error
	mutex   sync.Mutex
}

func (t *testRunnable) Run(context.Context) error {
	t.mutex.Lock()
	t.counter++
	t.mutex.Unlock()
	return t.result
}

func TestAllRunNothing(t *testing.T) {
	err := All(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestAllRunOne(t *testing.T) {
	r1 := new(testRunnable)
	err := All(context.Background(), r1.Run)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r1.counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestAllWithError(t *testing.T) {
	r1 := new(testRunnable)
	r1.result = errors.New("fail")
	r2 := new(testRunnable)
	err := All(context.Background(), r1.Run, r2.Run)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r1.counter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r2.counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestAllRunThree(t *testing.T) {
	r1 := new(testRunnable)
	err := All(context.Background(), r1.Run, r1.Run, r1.Run)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(r1.counter, Ge(1)); err != nil {
		t.Fatal(err)
	}
}
