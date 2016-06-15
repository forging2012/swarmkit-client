package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/shenshouer/swarmkit-client/api"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           os.Args[0],
	Short:         "The http client for swarmkit.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO add tls support
		var tlsConfig *tls.Config = nil

		enableCors, err := cmd.Flags().GetBool("api-enable-cors")
		if err != nil {
			log.Fatal(err)
		}
		host, err := cmd.Flags().GetString("advertise")
		if err != nil {
			log.Fatal(err)
		}
		socker, err := cmd.Flags().GetString("socket")
		if err != nil {
			log.Fatal(err)
		}
		server := api.NewServer(host, tlsConfig)
		swarmkitAPI, err := api.Dial(socker)
		if err != nil {
			log.Fatal(err)
		}
		primary := api.NewPrimary(swarmkitAPI, tlsConfig, enableCors)
		server.SetHandler(primary)
		log.Fatal(server.ListenAndServe())
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func defaultSocket() string {
	swarmSocket := os.Getenv("SWARM_SOCKET")
	if swarmSocket != "" {
		return swarmSocket
	}
	return "/var/run/docker/cluster/docker-swarmd.sock"
}

func init() {
	RootCmd.PersistentFlags().StringP("socket", "s", defaultSocket(), "Socket to connect to the Swarm manager")
	RootCmd.PersistentFlags().BoolP("api-enable-cors", "c", false, "enable CORS headers in the remote API (default false)")
	RootCmd.PersistentFlags().StringP("advertise", "a", ":8888", "advertise for http server")
}
