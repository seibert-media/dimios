package k8s_deploy

import (
	"context"

	"github.com/bborbe/k8s_deploy/apply"
	"github.com/bborbe/k8s_deploy/find"
	"github.com/bborbe/k8s_deploy/sync"
)

type Deployer struct {
	Dir string
}

func (d *Deployer) Deploy(ctx context.Context) error {
	changeFinder := find.NewFinder(find.ManifestDirectory(d.Dir))
	changeApplier := apply.NewApplier()
	changeSyncer := sync.NewSyncer(
		changeFinder.Changes,
		changeApplier.Apply,
	)
	return changeSyncer.Sync(ctx)
}
