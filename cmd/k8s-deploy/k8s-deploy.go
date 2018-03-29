package main

import (
	"context"
	"flag"
	"runtime"

	"github.com/golang/glog"
	"github.com/bborbe/k8s_deploy"
)

var dirPtr = flag.String("dir", "", "Path to manifest folder")

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	deployer := &k8s_deploy.Deployer{
		*dirPtr,
	}

	glog.V(0).Infof("k8s deploy started")
	defer glog.V(0).Infof("k8s deploy finished")
	if err := deployer.Deploy(context.Background()); err != nil {
		glog.Exit(err)
	}
}
