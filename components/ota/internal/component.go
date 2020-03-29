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
	logger logging.Logger
	routes []dashboard.Route

	installer        *ota.Installer
	uploadRepository *repository.Directory
	allRepository    *repository.Merge
	currentRelease   ota.Release
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

	c.installer = ota.NewInstaller(a.Shutdown)

	c.uploadRepository = repository.NewDirectory()
	c.allRepository = repository.NewMerge(c.uploadRepository, repository.NewMemory(c.currentRelease))

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) (err error) {
	c.logger = logging.DefaultLazyLogger(c.Name())

	<-a.ReadyComponent(config.ComponentName)
	cfg := a.GetComponent(config.ComponentName).(config.Component)

	c.uploadRepository.SetPath(cfg.String(ota.ConfigReleasesDirectory))
	go c.Update()

	shadowURLs := cfg.String(ota.ConfigRepositoryClientShadow)
	if shadowURLs != "" {
		for _, u := range strings.Split(shadowURLs, ",") {
			shadowURL, err := url.Parse(u)
			if err != nil {
				return err
			}

			c.allRepository.Merge(repository.NewShadow(shadowURL))
		}
	}

	return err
}

func (c *Component) Update() (err error) {
	err = c.allRepository.Update()
	if err != nil {
		c.logger.Error("Update all repositories failed", "error", err.Error())
	}

	return err
}
