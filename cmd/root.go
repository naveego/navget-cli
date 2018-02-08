// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "navget-cli [command]",
	Short: "Tool for creating and publishing navget packages.",
	Long:  `There must be a manifest.json file in the directory where you run this tool.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	if v, ok := os.LookupEnv("DRONE"); ok && v == "true" {
		log.Println("Drone host detected, using PLUGIN prefix rather than NAVGET for environment variables")
		viper.SetEnvPrefix("PLUGIN")
	} else {
		viper.SetEnvPrefix("NAVGET")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	RootCmd.PersistentFlags().StringP("package", "p", "package.zip", "The file name for package to create/upload.")
	RootCmd.PersistentFlags().StringP("files", "f", "", "The files to include in the package as a quoted, space-delimited list (like \"file1.exe file2.exe\").")
	viper.BindPFlags(RootCmd.Flags())

}
