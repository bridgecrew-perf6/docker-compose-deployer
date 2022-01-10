package cmd

import (
	"fmt"
	"net/http"

	"github.com/fatindeed/docker-compose-deployer/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	serverCmd.PersistentFlags().Int("port", 8000, "The server port")
	serverCmd.PersistentFlags().String("secret", "", "JWT secret")
	addDeployerFlags(serverCmd.PersistentFlags())
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		dh, err := services.NewDeployHandler()
		if err != nil {
			return err
		}
		http.Handle("/v1/deploy", dh)
		http.HandleFunc("/v1/ping", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "pong")
		})
		addr := fmt.Sprintf(":%d", viper.GetInt("port"))
		logrus.Infof("starting http server on %s", addr)
		return http.ListenAndServe(addr, nil)
	},
}
