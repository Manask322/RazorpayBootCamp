/*
Copyright © 2020 Manas Kinkar <manask322@gmai.com>

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
	"log"
	"os"

	"../db"
	"github.com/spf13/cobra"
)

var taskName string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:              "add",
	TraverseChildren: true,
	Short:            "add a taks to the list",

	Run: func(cmd *cobra.Command, args []string) {
		err := db.CreateRecord(taskName)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("\nTask Added")
	},
}

func init() {
	addCmd.Flags().StringVarP(&taskName, "task", "t", "", "Enter the task to add")
	addCmd.MarkFlagRequired("task")
	rootCmd.AddCommand(addCmd)
}
