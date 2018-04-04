package change

import (
	"fmt"

	"github.com/bborbe/k8s_deploy/k8s"
	"k8s.io/apimachinery/pkg/runtime"
)

type Change struct {
	Deleted bool
	Object  runtime.Object
}

func (c *Change) String() string {
	if c.Deleted {
		return fmt.Sprintf("DELETE %s", k8s.ObjectToString(c.Object))
	} else {
		return fmt.Sprintf("CREATE %s", k8s.ObjectToString(c.Object))
	}
}
