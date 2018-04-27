// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package change

import (
	"context"

	"github.com/bborbe/run"
	"github.com/golang/glog"
)

const channelSize = 10

//go:generate counterfeiter -o ../mocks/applier.go --fake-name Applier . applier
type applier interface {
	Run(context.Context, <-chan Change) error
}

//go:generate counterfeiter -o ../mocks/getter.go --fake-name Getter . getter
type getter interface {
	Run(context.Context, chan<- Change) error
}

type Syncer struct {
	Applier applier
	Getter  getter
}

// Run the sync until one function errors
func (s *Syncer) Run(ctx context.Context) error {
	glog.V(1).Info("sync changes started")
	defer glog.V(1).Info("sync changes finished")
	changes := make(chan Change, channelSize)
	defer close(changes)

	return run.CancelOnFirstError(ctx,
		// get changes
		func(ctx context.Context) error {
			return s.Getter.Run(ctx, changes)
		},
		// apply changes
		func(ctx context.Context) error {
			return s.Applier.Run(ctx, changes)
		},
	)
}
