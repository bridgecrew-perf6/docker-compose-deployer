package models

import (
	"regexp"
	"time"
)

// ComposeFile is a filename and the config data with last modified time
type ComposeFile struct {
	filename     string
	lastModified time.Time
	config       ComposeConfig
}

// ComposeConfig is a full compose file configuration and model
type ComposeConfig struct {
	Services map[string]ServiceConfig
}

func (c *ComposeConfig) merge(services map[string]ServiceConfig) {
	re := regexp.MustCompile(`(\S+):\${?([^}]+)}?`)
	for name, sc := range services {
		if sc.Image != "" {
			match := re.FindStringSubmatch(sc.Image)
			if match != nil {
				sc.ImageName = match[1]
				sc.TagName = match[2]
			}
			c.Services[name] = sc
		}
	}
}

// ServiceConfig is the configuration of one service
type ServiceConfig struct {
	Image     string `yaml:",omitempty"`
	ImageName string `yaml:"-"`
	TagName   string `yaml:"-"`
}
