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
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [files]",
	Short: "Creates a zip file containing the manifest.json and the provided files.",

	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())

		fmt.Println("create", viper.AllSettings())

		if len(args) == 0 {
			args = viper.GetStringSlice("files")
		}

		ExecuteCreate(args)
	},
}

func ExecuteCreate(files []string) {

	if _, err := os.Stat("manifest.json"); os.IsNotExist(err) {
		log.Fatalf("manifest.json file not found: %s", err)
	}

	outfile := "package.zip"

	log.Printf("Creating package at '%s'...", outfile)

	files = append(files, "manifest.json")

	written := map[string]bool{}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	for _, file := range files {
		// don't write files twice
		if _, alreadyWritten := written[file]; alreadyWritten {
			continue
		}
		written[file] = true
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("error reading file '%s': %s", file, err)
		}
		f, err := w.Create(file)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(outfile, buf.Bytes(), 0644)

	log.Printf("Created package at '%s'.", outfile)
}

func init() {
	RootCmd.AddCommand(createCmd)
	initCreateFlags(createCmd)

}

func initCreateFlags(cmd *cobra.Command) {
}
