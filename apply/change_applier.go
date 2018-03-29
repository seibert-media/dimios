package apply

import (
	"context"
	"fmt"

	"github.com/bborbe/k8s_deploy/change"
	"github.com/golang/glog"
)

type Applier struct {
}

func NewApplier() *Applier {
	return &Applier{}
}

func (c *Applier) Apply(ctx context.Context, changes <-chan change.Change) error {
	for {
		select {
		case v, ok := <-changes:
			if !ok {
				glog.V(1).Infoln("all changes applied")
				return nil
			}
			glog.V(3).Infof("apply change %v", v)
			if err := c.apply(v); err != nil {
				return fmt.Errorf("apply change failed: %v", err)
			}
		case <-ctx.Done():
			glog.V(3).Infoln("context done, skip apply changes")
			return nil
		}
	}
	return nil
}

func (c *Applier) apply(change change.Change) error {
	fmt.Printf("%v\n", change)
	return nil
}
