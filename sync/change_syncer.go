// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"

	"github.com/bborbe/run"
	"github.com/golang/glog"
	"github.com/seibert-media/k8s-deploy/change"
)

const channelSize = 10

// Syncer is responsible for sending incoming changes to the apply function
type Syncer interface {
	Run(ctx context.Context) error
}

// Handler interface for getting and applying changes
type Handler interface {
	Run(ctx context.Context, c chan change.Change) error
}

type syncer struct {
	getter  Handler
	applier Handler
}

// New Syncer taking get and apply functions
func New(
	getter Handler,
	applier Handler,
) Syncer {
	return &syncer{
		getter:  getter,
		applier: applier,
	}
}

// Run the sync until one function errors
func (c *syncer) Run(ctx context.Context) error {
	glog.V(1).Info("sync changes started")
	defer glog.V(1).Info("sync changes finished")
	versionChannel := make(chan change.Change, channelSize)

	return run.CancelOnFirstError(ctx,
		// get changes
		func(ctx context.Context) error {
			return c.getter.Run(ctx, versionChannel)
		},
		// apply changes
		func(ctx context.Context) error {
			return c.applier.Run(ctx, versionChannel)
		},
	)
}
