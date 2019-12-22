package internal

import (
	"net/url"
	"strings"

	"github.com/kardianos/osext"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/release"
	"github.com/kihamo/shadow/components/ota/repository"
)

type Component struct {
	config config.Component
	routes []dashboard.Route

	updater           *ota.Updater
	uploadRepository  *repository.Directory
	upgradeRepository *repository.Merge
	currentRelease    ota.Release
}

func (c *Component) Name() string {
	return ota.ComponentName
}

func (c *Component) Version() string {
	return ota.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: i18n.ComponentName,
		},
		{
			Name: logging.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	releasePath, err := osext.Executable()
	if err != nil {
		return err
	}

	c.currentRelease, err = release.NewLocalFile(releasePath, a.Version()+" "+a.Build())
	if err != nil {
		return err
	}

	c.updater = ota.NewUpdater()

	c.uploadRepository = repository.NewDirectory()
	c.uploadRepository.Add(c.currentRelease)

	c.upgradeRepository = repository.NewMerge(c.uploadRepository)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) (err error) {
	<-a.ReadyComponent(config.ComponentName)
	cfg := a.GetComponent(config.ComponentName).(config.Component)

	err = c.uploadRepository.Load(cfg.String(ota.ConfigReleasesDirectory))
	if err != nil {
		return err
	}

	shadowURLs := cfg.String(ota.ConfigRepositoryClientShadow)
	if shadowURLs != "" {
		for _, u := range strings.Split(shadowURLs, ",") {
			shadowURL, err := url.Parse(u)
			if err != nil {
				return err
			}

			c.upgradeRepository.Merge(repository.NewShadow(shadowURL))
		}
	}

	return err
}

func (c *Component) doAutoUpgrade() {
	// TODO:
}
