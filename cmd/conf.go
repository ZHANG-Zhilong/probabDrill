/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Path string = ""

// confCmd represents the conf command
var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "setup for config file.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("conf called")
		fmt.Println(args)
		fmt.Println(Path)
	},
}

func init() {
	rootCmd.AddCommand(confCmd)
	confCmd.Flags().StringVarP(&Path, "path", "p", "./conf", "path of config")
}
