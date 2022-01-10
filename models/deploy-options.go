package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	defaultFileNames         = []string{"compose.yaml", "compose.yml", "docker-compose.yml", "docker-compose.yaml"}
	defaultOverrideFileNames = []string{"compose.override.yml", "compose.override.yaml", "docker-compose.override.yml", "docker-compose.override.yaml"}
)

type DeployOptions struct {
	// docker compose arguments
	ProjectName string   `mapstructure:"project-name"`
	ConfigPaths []string `mapstructure:"file"`
	WorkDir     string   `mapstructure:"workdir"`
	// runtime arguments
	ComposeV2 bool `mapstructure:"compose-v2"`
	Sudo      bool `mapstructure:"sudo"`
	// runtime variables
	composeFiles  []*ComposeFile
	composeConfig *ComposeConfig
}

func (o *DeployOptions) Init() error {
	stat, err := os.Stat(o.WorkDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", o.WorkDir)
	}
	_, err = os.Stat(o.GetDotEnv())
	if err != nil {
		return fmt.Errorf("stat .env error: %w", err)
	}

	o.composeFiles = make([]*ComposeFile, 0)
	if len(o.ConfigPaths) == 0 {
		candidates := findFiles(defaultFileNames, o.WorkDir)
		if len(candidates) == 0 {
			return fmt.Errorf("no compose file found")
		}
		o.composeFiles = append(o.composeFiles, &ComposeFile{filename: candidates[0]})
		overrides := findFiles(defaultOverrideFileNames, o.WorkDir)
		if len(overrides) > 0 {
			o.composeFiles = append(o.composeFiles, &ComposeFile{filename: candidates[0]})
		}
	} else {
		for _, filename := range o.ConfigPaths {
			compose := &ComposeFile{filename: filename}
			o.composeFiles = append(o.composeFiles, compose)
		}
	}
	return o.LoadComposeConfig()
}

func (o *DeployOptions) GetDotEnv() string {
	return filepath.Join(o.WorkDir, ".env")
}

func (o *DeployOptions) GetComposeArgs() []string {
	args := make([]string, 0)
	if o.ProjectName != "" {
		args = append(args, "--project-name", o.ProjectName)
	}
	if len(o.ConfigPaths) > 0 {
		for _, file := range o.ConfigPaths {
			args = append(args, "-f", file)
		}
	}
	if o.WorkDir != "" {
		if o.ComposeV2 {
			args = append(args, "--project-directory", o.WorkDir)
		} else {
			args = append(args, "--workdir", o.WorkDir)
		}
	}
	return args
}

func (o *DeployOptions) LoadComposeConfig() error {
	var modified bool
	composeConfig := new(ComposeConfig)
	composeConfig.Services = make(map[string]ServiceConfig)
	for _, compose := range o.composeFiles {
		stat, err := os.Stat(compose.filename)
		if err != nil {
			return fmt.Errorf("stat %s error: %w", compose.filename, err)
		}
		lastModified := stat.ModTime()
		if compose.lastModified.Equal(lastModified) {
			continue
		}
		modified = true
		data, err := os.ReadFile(compose.filename)
		if err != nil {
			return fmt.Errorf("error loading %s: %w", compose.filename, err)
		}
		err = yaml.Unmarshal(data, &compose.config)
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", compose.filename, err)
		}
		composeConfig.merge(compose.config.Services)
		compose.lastModified = lastModified
	}
	if modified {
		o.composeConfig = composeConfig
	}
	return nil
}

func (o *DeployOptions) GetServiceConfig(svc string) (*ServiceConfig, error) {
	sc, ok := o.composeConfig.Services[svc]
	if !ok {
		return nil, fmt.Errorf("service %s not found", svc)
	}
	if sc.ImageName == "" || sc.TagName == "" {
		return nil, fmt.Errorf("invalid image format: %s", sc.Image)
	}
	return &sc, nil
}

func findFiles(names []string, pwd string) []string {
	candidates := []string{}
	for _, n := range names {
		f := filepath.Join(pwd, n)
		if _, err := os.Stat(f); err == nil {
			candidates = append(candidates, f)
		}
	}
	if len(candidates) > 1 {
		logrus.Warnf("more than one file found")
	}
	return candidates
}
