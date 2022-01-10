package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("DEPLOYER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

// Version is the version of the CLI injected in compilation time
var (
	verbose bool
	Version = "dev"
	rootCmd = &cobra.Command{
		Use:     "docker-compose-deployer",
		Version: fmt.Sprintf("1.0.0, build %s", Version),
		Short:   "Deployer for your Docker Compose services",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlags(cmd.PersistentFlags())
			if err != nil {
				logrus.Error(err)
			}
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		SilenceUsage: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addDeployerFlags(f *pflag.FlagSet) {
	// docker compose configuration
	f.StringP("project-name", "p", "", "Compose project name")
	f.StringArrayP("file", "f", []string{}, "Compose configuration files")
	f.StringP("workdir", "w", ".", "Specify an alternate working directory")
	// runtime configuration
	f.Bool("compose-v2", false, "Use Docker Compose v2")
	f.Bool("sudo", false, "Execute commands with sudo")
	f.BoolVarP(&verbose, "verbose", "v", false, "Show more output")
}
