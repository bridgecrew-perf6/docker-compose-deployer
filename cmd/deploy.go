package cmd

import (
	"github.com/fatindeed/docker-compose-deployer/services"
	"github.com/spf13/cobra"
)

func init() {
	addDeployerFlags(deployCmd.PersistentFlags())
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy <service name> <image tag>",
	Short: "Deploy service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deployer, err := services.NewDeployer()
		if err != nil {
			return err
		}
		return deployer.Run(args[0], args[1])
	},
}
