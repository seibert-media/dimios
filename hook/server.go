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
	glog.V(1).Info("sync changes triggerd")
	if err := s.Manager.Run(req.Context()); err != nil {
		glog.Warningf("sync changes failed: %v", err)
		http.Error(resp, fmt.Sprintf("run failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	glog.V(1).Info("sync changes completed successful")
	fmt.Fprintln(resp, "ok")
}
