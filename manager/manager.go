// Copyright 2018 The K8s-Deploy Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/bborbe/http/client_builder"
	"github.com/bborbe/teamvault_utils/connector"
	"github.com/bborbe/teamvault_utils/model"
	"github.com/bborbe/teamvault_utils/parser"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/k8s-deploy/apply"
	"github.com/seibert-media/k8s-deploy/finder"
	"github.com/seibert-media/k8s-deploy/k8s"
	file_provider "github.com/seibert-media/k8s-deploy/k8s/file"
	remote_provider "github.com/seibert-media/k8s-deploy/k8s/remote"
	"github.com/seibert-media/k8s-deploy/sync"
	k8s_discovery "k8s.io/client-go/discovery"
	k8s_dynamic "k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Required for using GCP auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// Manager is the main application package
type Manager struct {
	Staging             bool
	TemplateDirectory   string
	TeamvaultURL        string
	TeamvaultUser       string
	TeamvaultPassword   string
	TeamvaultConfigPath string
	Namespaces          string
	Kubeconfig          string
}

// ReadTeamvaultConfig from path
func (m *Manager) ReadTeamvaultConfig() error {
	teamvaultConfigPath := model.TeamvaultConfigPath(m.TeamvaultConfigPath)
	if teamvaultConfigPath.Exists() {
		teamvaultConfig, err := teamvaultConfigPath.Parse()
		if err != nil {
			glog.V(2).Infof("parse teamvault config failed: %v", err)
			return err
		}
		m.TeamvaultURL = teamvaultConfig.Url.String()
		m.TeamvaultUser = teamvaultConfig.User.String()
		m.TeamvaultPassword = teamvaultConfig.Password.String()
	}
	return nil
}

// Validate if all Manager values are set correctly
func (m *Manager) Validate() error {
	if len(m.TemplateDirectory) == 0 {
		return fmt.Errorf("template directory missing")
	}
	if len(m.Namespaces) == 0 {
		return fmt.Errorf("namespace missing")
	}
	if len(m.Kubeconfig) == 0 {
		return fmt.Errorf("kubeconfig missing")
	}
	if len(m.TeamvaultURL) == 0 && !m.Staging {
		return fmt.Errorf("teamvault url missing")
	}
	if len(m.TeamvaultUser) == 0 && !m.Staging {
		return fmt.Errorf("teamvault user missing")
	}
	if len(m.TeamvaultPassword) == 0 && !m.Staging {
		return fmt.Errorf("teamvault password missing")
	}
	return nil
}

// Run Manager
func (m *Manager) Run(ctx context.Context) error {
	glog.V(0).Info("kubernetes-manager started")
	defer glog.V(0).Info("kubernetes-manager finished")

	discovery, dynamicPool, err := m.createClients()
	if err != nil {
		return err
	}

	fileProvider := file_provider.New(
		file_provider.TemplateDirectory(m.TemplateDirectory),
		m.createTeamvaultConfigParser(),
	)
	remoteProvider := remote_provider.New(discovery, dynamicPool)

	syncer := sync.New(
		finder.New(
			fileProvider,
			remoteProvider,
			k8s.NamespacesFromCommaSeperatedList(m.Namespaces),
		),
		apply.New(
			m.Staging,
			discovery,
			dynamicPool,
		),
	)

	return syncer.Run(ctx)
}

func (m *Manager) createClientConfig() (*restclient.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", m.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("build config from flags failed: %v", err)
	}
	return config, nil
}

func (m *Manager) createClients() (*k8s_discovery.DiscoveryClient, k8s_dynamic.ClientPool, error) {
	cfg, err := m.createClientConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "create clientConfig failed")
	}

	discovery, err := k8s_discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "creating k8s_discovery client failed")
	}
	dynamicPool := k8s_dynamic.NewDynamicClientPool(cfg)

	return discovery, dynamicPool, nil
}

func (m *Manager) createTeamvaultConfigParser() parser.Parser {
	return parser.New(m.createTeamvaultConnector())
}

func (m *Manager) createTeamvaultConnector() connector.Connector {
	var teamvaultConnector connector.Connector
	if m.Staging {
		teamvaultConnector = connector.NewDummy()
	}
	httpClient := client_builder.New().WithTimeout(5 * time.Second).Build()
	teamvaultConnector = connector.New(
		httpClient.Do,
		model.TeamvaultUrl(m.TeamvaultURL),
		model.TeamvaultUser(m.TeamvaultUser),
		model.TeamvaultPassword(m.TeamvaultPassword),
	)
	return teamvaultConnector
}
