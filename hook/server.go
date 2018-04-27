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

type Server struct {
	Manager manager
}

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(1).Info("sync changes for started")
	if err := s.Manager.Run(req.Context()); err != nil {
		glog.V(0).Info("sync changes failed: %v", err)
		http.Error(resp, fmt.Sprintf("run failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	glog.V(1).Info("sync changes finished")
	fmt.Fprintln(resp, "ok")
}
