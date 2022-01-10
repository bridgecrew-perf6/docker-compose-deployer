package services

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"

	"github.com/fatindeed/docker-compose-deployer/models"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type deployer struct {
	options  *models.DeployOptions
	locker   sync.Mutex
	commands map[string][]string
}

func (d *deployer) Run(svc, tag string) error {
	d.locker.Lock()
	defer d.locker.Unlock()

	logrus.Debugf("deploying %s with tag %#v", svc, tag)

	// Check if docker-compose.yml been modified
	err := d.options.LoadComposeConfig()
	if err != nil {
		return err
	}

	// Get service image info from config
	sc, err := d.options.GetServiceConfig(svc)
	if err != nil {
		return err
	}

	// Test to pull new image
	image := fmt.Sprintf("%s:%s", sc.ImageName, tag)
	logrus.Debugf("pulling image %s", image)
	out, err := d.exec("docker", "pull", image)
	if err != nil {
		return fmt.Errorf("error pulling %s: %w\n%s", image, err, string(out))
	}

	// Update .env file
	err = d.updateEnv(sc.TagName, tag)
	if err != nil {
		return err
	}

	// Restart docker compose service
	logrus.Debugf("restarting service %s", svc)
	out, err = d.exec("docker-compose", append(d.options.GetComposeArgs(), "up", "-d", "--remove-orphans", svc)...)
	if err != nil {
		return fmt.Errorf("error restarting %s: %w\n%s", svc, err, string(out))
	}

	logrus.Infof("%s deployed with tag: %s", svc, tag)
	return nil
}

func (d *deployer) updateEnv(key, value string) error {
	dotenv := d.options.GetDotEnv()
	envMap, err := godotenv.Read(dotenv)
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	envMap[key] = value
	godotenv.Write(envMap, dotenv)
	return nil
}

func (d *deployer) exec(name string, args ...string) ([]byte, error) {
	logrus.Debugf("%s: %v", name, d.commands[name])
	args = append(d.commands[name], args...)
	if d.options.Sudo {
		args = append([]string{"sudo"}, args...)
	}
	logrus.Debugf("executing command: %v", args)
	outbuf := bytes.NewBuffer(nil)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = outbuf
	cmd.Stderr = outbuf
	cmd.Dir = d.options.WorkDir
	err := cmd.Run()
	logrus.Debugf("result: %s", outbuf.String())
	return outbuf.Bytes(), err
}

func NewDeployer() (*deployer, error) {
	d := new(deployer)
	err := viper.Unmarshal(&d.options)
	if err != nil {
		return nil, err
	}
	err = d.options.Init()
	if err != nil {
		return nil, err
	}
	logrus.Infof("starting deployer with %#v", d.options)

	d.commands = make(map[string][]string)
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return nil, err
	}
	d.commands["docker"] = []string{dockerPath}
	if d.options.ComposeV2 {
		d.commands["docker-compose"] = []string{dockerPath, "compose"}
	} else {
		composePath, err := exec.LookPath("docker-compose")
		if err != nil {
			return nil, err
		}
		d.commands["docker-compose"] = []string{composePath}
	}

	return d, nil
}
