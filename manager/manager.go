// Copyright 2018 The Dimios Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bborbe/http/client_builder"
	"github.com/bborbe/teamvault_utils/connector"
	"github.com/bborbe/teamvault_utils/model"
	"github.com/bborbe/teamvault_utils/parser"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/dimios/apply"
	"github.com/seibert-media/dimios/change"
	"github.com/seibert-media/dimios/finder"
	"github.com/seibert-media/dimios/k8s"
	file_provider "github.com/seibert-media/dimios/k8s/file"
	remote_provider "github.com/seibert-media/dimios/k8s/remote"
	k8s_discovery "k8s.io/client-go/discovery"
	k8s_dynamic "k8s.io/client-go/dynamic"
	k8s_rest "k8s.io/client-go/rest"
	k8s_clientcmd "k8s.io/client-go/tools/clientcmd"

	// Required for using GCP auth
	"os"

	"github.com/seibert-media/dimios/hook"
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
	Webhook             bool
	Port                int
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
	if len(m.Kubeconfig) == 0 && os.Getenv("KUBERNETES_SERVICE_HOST") == "" && os.Getenv("KUBERNETES_SERVICE_PORT") == "" {
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

	applier := apply.New(
		m.Staging,
		discovery,
		dynamicPool,
	)
	getter := finder.New(
		fileProvider,
		remoteProvider,
		k8s.NamespacesFromCommaSeperatedList(m.Namespaces),
	)
	syncer := &change.Syncer{
		Applier: applier,
		Getter:  getter,
	}

	if m.Webhook {
		server := &http.Server{
			Addr: fmt.Sprintf(":%d", m.Port),
			Handler: &hook.Server{
				Manager: syncer,
			},
		}
		glog.V(1).Infof("start webserver on port %d", m.Port)
		return server.ListenAndServe()
	}
	glog.V(1).Infof("run sync")
	return syncer.Run(ctx)
}

func createConfig(kubeconfig string) (*k8s_rest.Config, error) {
	if len(kubeconfig) > 0 {
		glog.V(4).Infof("create kube config from flags")
		return k8s_clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	glog.V(4).Infof("create in cluster kube config")
	return k8s_rest.InClusterConfig()
}

func (m *Manager) createClients() (*k8s_discovery.DiscoveryClient, k8s_dynamic.ClientPool, error) {
	cfg, err := createConfig(m.Kubeconfig)
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
