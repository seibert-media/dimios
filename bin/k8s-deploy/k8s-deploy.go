package main

import (
	"context"
	"flag"
	"runtime"

	"github.com/bborbe/k8s_deploy/apply"
	"github.com/bborbe/k8s_deploy/find"
	"github.com/bborbe/k8s_deploy/sync"
	"github.com/golang/glog"
)

var dirPtr = flag.String("dir", "", "Path to fanifest folder")

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := do(context.Background()); err != nil {
		glog.Exit(err)
	}
}

func do(ctx context.Context) error {
	glog.V(0).Infof("k8s deploy started")
	defer glog.V(0).Infof("k8s deploy finished")
	syncer := createSyncer()
	return syncer.Sync(ctx)
}

func createSyncer() sync.Syncer {
	changeFinder := find.NewFinder(find.ManifestDirectory(*dirPtr))
	changeApplier := apply.NewApplier()
	changeSyncer := sync.NewSyncer(
		changeFinder.Changes,
		changeApplier.Apply,
	)
	return changeSyncer
}
