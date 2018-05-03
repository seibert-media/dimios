package hook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

//go:generate counterfeiter -o ../mocks/manager.go --fake-name Manager . manager
type manager interface {
	Run(ctx context.Context) error
}

type handler struct {
	Manager manager
	running chan bool
}

func NewHandler(manager manager) *handler {
	return &handler{
		Manager: manager,
		running: make(chan bool, 1),
	}
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(1).Info("sync changes triggerd")
	select {
	case h.running <- true:
		go func() {
			defer func() { <-h.running }()
			if err := h.Manager.Run(context.Background()); err != nil {
				glog.Warningf("sync changes failed: %v", err)
			}
			glog.V(1).Info("sync changes completed successful")
		}()
		fmt.Fprintln(resp, "sync triggerd")
	default:
		fmt.Fprintln(resp, "already running")
	}
}
