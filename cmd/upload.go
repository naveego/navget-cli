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
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the publish command
var uploadCmd = &cobra.Command{
	Use:   "upload [options]",
	Short: "Uploads a package to Navget.",
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())

		fmt.Println("upload", viper.AllSettings())

		os := requireParam("os")
		arch := requireParam("arch")

		ExecutePublish(os, arch)
	},
}

func ExecutePublish(osName, arch string) {
	packageFilePath := requireParam("package")
	log.Printf("Uploading '%s'...", packageFilePath)

	endpoint := requireParam("endpoint")
	url := fmt.Sprintf("%s/api/packages?os=%s&arch=%s", endpoint, osName, arch)
	log.Printf("Uploading to '%s'...", url)

	token := requireParam("token")

	file, err := os.Open(packageFilePath)
	defer file.Close()
	if err != nil {
		log.Fatalf("error reading package file: %s", err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(packageFilePath))
	if err != nil {
		log.Fatal("couldn't create form file: ", err)
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		log.Fatal("couldn't close form file: ", err)
	}
	log.Printf("Upload size: %d bytes", body.Len())
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal("couldn't create request: ", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("couldn't upload file: ", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	json := buf.String()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error uploading - status=%v, response=%s", resp.StatusCode, json)
	}

	defer resp.Body.Close()
	log.Printf("Uploaded '%s' to '%s'.", packageFilePath, url)

}

func requireParam(name string) string {

	p := viper.GetString(name)
	if p == "" {
		log.Fatalf("parameter '%s' (or env NAVGET_%s) is required", name, strings.Replace(strings.ToUpper(name), "-", "_", -1))
	}
	return p
}

func initUploadFlags(cmd *cobra.Command) {
	cmd.Flags().String("os", "linux", "The OS the package to.")
	cmd.Flags().String("arch", "amd64", "The architecture for the package to.")
	cmd.Flags().StringP("endpoint", "e", "", "The endpoint to upload the package to.")
	cmd.Flags().StringP("token", "t", "", "The JWT token used to authenticate to the Navget endpoint.")
	viper.BindEnv("token", "NAVGET_TOKEN", "PLUGIN_TOKEN")
}

func init() {

	initUploadFlags(uploadCmd)
	RootCmd.AddCommand(uploadCmd)

}
