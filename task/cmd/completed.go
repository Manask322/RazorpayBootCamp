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
	"log"
	"os"

	"../db"
	"github.com/spf13/cobra"
)

// completedCmd represents the completed command
var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "lists all the completed tasks",

	Run: func(cmd *cobra.Command, args []string) {
		err := db.CompletedTasks()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
