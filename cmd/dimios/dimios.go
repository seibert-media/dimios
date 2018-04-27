// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/golang/glog"
	"github.com/kolide/kit/version"
	"github.com/seibert-media/dimios/manager"
)

var (
	templateDirectoryPtr   = flag.String("dir", "", "Path to template directory")
	namespacesPtr          = flag.String("namespaces", "", "list of kubernetes namespace separated by comma")
	teamvaultURLPtr        = flag.String("teamvault-url", "", "teamvault url")
	teamvaultUserPtr       = flag.String("teamvault-user", "", "teamvault user")
	teamvaultPassPtr       = flag.String("teamvault-pass", "", "teamvault password")
	teamvaultConfigPathPtr = flag.String("teamvault-config", "", "teamvault config")
	stagingPtr             = flag.Bool("staging", false, "staging status")
	versionInfo            = flag.Bool("version", false, "show version info")
	kubeconfig             = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	port                   = flag.Int("port", 8080, "port listen on if webhook is activated")
	webhook                = flag.Bool("webhook", false, "activate run as http server")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *versionInfo {
		fmt.Printf("-- Dimios --\n")
		version.PrintFull()
		os.Exit(0)
	}

	m := &manager.Manager{
		Namespaces:          *namespacesPtr,
		TemplateDirectory:   *templateDirectoryPtr,
		Staging:             *stagingPtr,
		TeamvaultConfigPath: *teamvaultConfigPathPtr,
		TeamvaultURL:        *teamvaultURLPtr,
		TeamvaultUser:       *teamvaultUserPtr,
		TeamvaultPassword:   *teamvaultPassPtr,
		Kubeconfig:          *kubeconfig,
		Webhook:             *webhook,
		Port:                *port,
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
