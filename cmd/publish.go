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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Creates and uploads a package.",
	Long:  `This invokes both the create and the upload commands.`,
	Run: func(cmd *cobra.Command, args []string) {

		viper.BindPFlags(cmd.Flags())

		log.Println("publish settings: ", viper.AllSettings())

		log.Println("ENV: ", os.Environ())

		if len(args) == 0 {
			args = viper.GetStringSlice("files")
		}

		ExecuteCreate(args)

		os := requireParam("os")
		arch := requireParam("arch")
		ExecutePublish(os, arch)
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)

	initCreateFlags(publishCmd)
	initUploadFlags(publishCmd)
}
