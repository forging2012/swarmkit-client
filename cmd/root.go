// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           os.Args[0],
	Short:         "The http client for swarmkit.",
	SilenceUsage:  true,
	SilenceErrors: true,
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
	// cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringP("socket", "s", defaultSocket(), "Socket to connect to the Swarm manager")
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.swarmkit-client.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" { // enable ability to specify config file via flag
// 		viper.SetConfigFile(cfgFile)
// 	}

// 	viper.SetConfigName(".swarmkit-client") // name of config file (without extension)
// 	viper.AddConfigPath("$HOME")            // adding home directory as first search path
// 	viper.AutomaticEnv()                    // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
