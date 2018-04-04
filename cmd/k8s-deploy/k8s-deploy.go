package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"path/filepath"

	"github.com/bborbe/k8s_deploy/manager"
	"github.com/golang/glog"
	"k8s.io/client-go/util/homedir"
)

var (
	templateDirectoryPtr   = flag.String("dir", "", "Path to template directory")
	namespacePtr           = flag.String("namespace", "", "Kubernetes namespace")
	teamvaultUrlPtr        = flag.String("teamvault-url", "", "teamvault url")
	teamvaultUserPtr       = flag.String("teamvault-user", "", "teamvault user")
	teamvaultPassPtr       = flag.String("teamvault-pass", "", "teamvault password")
	teamvaultConfigPathPtr = flag.String("teamvault-config", "", "teamvault config")
	stagingPtr             = flag.Bool("staging", false, "staging status")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	m := &manager.Manager{
		Namespace:           *namespacePtr,
		TemplateDirectory:   *templateDirectoryPtr,
		Staging:             *stagingPtr,
		TeamvaultConfigPath: *teamvaultConfigPathPtr,
		TeamvaultUrl:        *teamvaultUrlPtr,
		TeamvaultUser:       *teamvaultUserPtr,
		TeamvaultPassword:   *teamvaultPassPtr,
		Kubeconfig:          *kubeconfig,
	}

	if err := m.ReadTeamvaultConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "read teamvault config failed: %v\n", err.Error())
		os.Exit(1)
	}

	if err := m.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "parameter invalid: %v\n", err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := m.Run(context.Background()); err != nil {
		glog.Exit(err)
	}
}
