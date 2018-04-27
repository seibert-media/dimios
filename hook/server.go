package hook

import (
	"context"
	"fmt"
	"net/http"
)

//go:generate counterfeiter -o ../mocks/manager.go --fake-name Manager . manager
type manager interface {
	Run(ctx context.Context) error
}

type Server struct {
	Manager manager
}

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if err := s.Manager.Run(req.Context()); err != nil {
		http.Error(resp, fmt.Sprintf("run failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(resp, "ok")
}
