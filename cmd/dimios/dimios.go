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
	kubeconfig          = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	namespaces          = flag.String("namespaces", "", "list of kubernetes namespace separated by comma")
	port                = flag.Int("port", 8080, "port listen on if webhook is activated")
	staging             = flag.Bool("staging", false, "staging status")
	teamvaultConfigPath = flag.String("teamvault-config", "", "teamvault config")
	teamvaultPass       = flag.String("teamvault-pass", "", "teamvault password")
	teamvaultURL        = flag.String("teamvault-url", "", "teamvault url")
	teamvaultUser       = flag.String("teamvault-user", "", "teamvault user")
	templateDirectory   = flag.String("dir", "", "Path to template directory")
	versionInfo         = flag.Bool("version", false, "show version info")
	webhook             = flag.Bool("webhook", false, "activate run as http server")
	whitelistPtr        = flag.String("whitelist", "", "list of objecttypes separated by comma")
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

	glog.V(1).Infof("parameter Kubeconfig: %s", *kubeconfig)
	glog.V(1).Infof("parameter Namespaces: %s", *namespaces)
	glog.V(1).Infof("parameter Port: %s", *port)
	glog.V(1).Infof("parameter Staging: %s", *staging)
	glog.V(1).Infof("parameter TeamvaultConfigPath: %s", *teamvaultConfigPath)
	glog.V(1).Infof("parameter TeamvaultPassword: length=%d", len(*teamvaultPass))
	glog.V(1).Infof("parameter TeamvaultURL: %s", *teamvaultURL)
	glog.V(1).Infof("parameter TeamvaultUser: %s", *teamvaultUser)
	glog.V(1).Infof("parameter TemplateDirectory: %s", *templateDirectory)
	glog.V(1).Infof("parameter Webhook: %s", *webhook)
	glog.V(1).Infof("parameter Whitelist: %s", *whitelistPtr)

	m := &manager.Manager{
		Kubeconfig:          *kubeconfig,
		Namespaces:          *namespaces,
		Port:                *port,
		Staging:             *staging,
		TeamvaultConfigPath: *teamvaultConfigPath,
		TeamvaultPassword:   *teamvaultPass,
		TeamvaultURL:        *teamvaultURL,
		TeamvaultUser:       *teamvaultUser,
		TemplateDirectory:   *templateDirectory,
		Webhook:             *webhook,
		Whitelist:           *whitelistPtr,
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
