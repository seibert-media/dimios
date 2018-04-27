// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"

	"github.com/bborbe/run"
	"github.com/golang/glog"
	"github.com/seibert-media/dimios/change"
)

const channelSize = 10

type applier interface {
	Run(context.Context, chan<- change.Change) error
}

type getter interface {
	Run(context.Context, <-chan change.Change) error
}

// Run the sync until one function errors
func Run(ctx context.Context, applier applier, getter getter) error {
	glog.V(1).Info("sync changes started")
	defer glog.V(1).Info("sync changes finished")
	versionChannel := make(chan change.Change, channelSize)

	return run.CancelOnFirstError(ctx,
		// get changes
		func(ctx context.Context) error {
			return getter.Run(ctx, versionChannel)
		},
		// apply changes
		func(ctx context.Context) error {
			return applier.Run(ctx, versionChannel)
		},
	)
}
