// Copyright © 2017 Daniel Manesku <daniel.manesku@gmail.com>
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
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List items in the trash can",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-25s%-25s\n", "DELETION ID", "DELETION TIME")

		fmt.Println(time.Now())

		trashPath := getTrashPath()
		trashDirs, _ := ioutil.ReadDir(trashPath)
		for _, dir := range trashDirs {
			dirName := dir.Name()
			fmt.Println(dirName)
		}

	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
