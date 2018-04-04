package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/bborbe/http/client_builder"
	"github.com/bborbe/k8s_deploy/apply"
	"github.com/bborbe/k8s_deploy/finder"
	"github.com/bborbe/k8s_deploy/k8s"
	"github.com/bborbe/k8s_deploy/k8s/file"
	"github.com/bborbe/k8s_deploy/k8s/remote_provider"
	"github.com/bborbe/k8s_deploy/sync"
	"github.com/bborbe/teamvault_utils/connector"
	"github.com/bborbe/teamvault_utils/model"
	"github.com/bborbe/teamvault_utils/parser"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Manager struct {
	Staging             bool
	TemplateDirectory   string
	TeamvaultUrl        string
	TeamvaultUser       string
	TeamvaultPassword   string
	TeamvaultConfigPath string
	Namespace           string
	Kubeconfig          string
}

func (m *Manager) ReadTeamvaultConfig() error {
	teamvaultConfigPath := model.TeamvaultConfigPath(m.TeamvaultConfigPath)
	if teamvaultConfigPath.Exists() {
		teamvaultConfig, err := teamvaultConfigPath.Parse()
		if err != nil {
			glog.V(2).Infof("parse teamvault config failed: %v", err)
			return err
		}
		m.TeamvaultUrl = teamvaultConfig.Url.String()
		m.TeamvaultUser = teamvaultConfig.User.String()
		m.TeamvaultPassword = teamvaultConfig.Password.String()
	}
	return nil
}

func (m *Manager) Validate() error {
	if len(m.TemplateDirectory) == 0 {
		return fmt.Errorf("template directory missing")
	}
	if len(m.Namespace) == 0 {
		return fmt.Errorf("namespace missing")
	}
	if len(m.Kubeconfig) == 0 {
		return fmt.Errorf("kubeconfig missing")
	}
	if !m.Staging {
		if len(m.TeamvaultUrl) == 0 {
			return fmt.Errorf("teamvaul url missing")
		}
		if len(m.TeamvaultUser) == 0 {
			return fmt.Errorf("teamvaul user missing")
		}
		if len(m.TeamvaultPassword) == 0 {
			return fmt.Errorf("teamvaul password missing")
		}
	}
	return nil
}

func (m *Manager) Run(ctx context.Context) error {
	glog.V(0).Info("kubernetes-manager started")
	defer glog.V(0).Info("kubernetes-manager finished")

	fileProvider := file.New(file.TemplateDirectory(m.TemplateDirectory), m.createTeamvaultConfigParser())
	clientConfig, err := m.createClientConfig()
	if err != nil {
		return fmt.Errorf("create clientConfig failed: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return fmt.Errorf("create clientSet failed: %v", err)
	}
	removeProvider := remote_provider.New(clientSet)

	changeFinder := &finder.Finder{
		FileProvider:   fileProvider,
		RemoveProvider: removeProvider,
		Namespace:      k8s.Namespace(m.Namespace),
	}
	changeApplier, err := apply.New(clientConfig)
	if err != nil {
		return errors.Wrap(err, "creating applier failed")
	}
	changeSyncer := sync.New(
		changeFinder.Changes,
		changeApplier.Apply,
	)

	return changeSyncer.SyncChanges(ctx)
}

func (m *Manager) createClientConfig() (*restclient.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", m.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("build config from flags failed: %v", err)
	}
	return config, nil
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
		model.TeamvaultUrl(m.TeamvaultUrl),
		model.TeamvaultUser(m.TeamvaultUser),
		model.TeamvaultPassword(m.TeamvaultPassword),
	)
	return teamvaultConnector
}
