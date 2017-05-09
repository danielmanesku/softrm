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
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

const flush_all_flag string = "all"

// flushCmd represents the flush command
var flushCmd = &cobra.Command{
	Use:   "flush ID [ID...]",
	Short: "Delete files permanently",
	Long: `Delete given deletion groups permanently from the disk.
Use listing command to find deletion group IDs of files
from trash. Note that files that are flushed cannot be
recovered with this tool. Notice also that no file shredding
will be done - if you need shredding, use a specialized tool for it.`,
	Run: func(cmd *cobra.Command, args []string) {
		allFlag, _ := cmd.Flags().GetBool(flush_all_flag)

		// validate inputs
		{
			if !allFlag && len(args) < 1 { // not enough parameters
				cmd.Help()
				return
			}

			// if --all flag is passed, there should not be any other argument
			if allFlag && len(args) > 0 {
				fmt.Println("Error: --all parameter cannot be combined with other arguments")
				fmt.Println()
				cmd.Help()
				return
			}
		}

		ensureTrashDirExists(false)
		trashPath := getTrashPath()

		// handle the case when --all is passed
		if allFlag {
			os.RemoveAll(trashPath)
			createDir(trashPath)
			return
		}

		// mapping between all available deletion group ids and their directory names
		delGroupMap := make(map[string]string)

		// check validity of flags
		{
			files, _ := ioutil.ReadDir(trashPath)
			for _, f := range files {
				fileName := f.Name()
				delGroupId := fileName[strings.LastIndex(fileName, "-")+1:]
				delGroupMap[delGroupId] = fileName
			}
			for _, arg := range args {
				_, found := delGroupMap[arg]
				if !found {
					fmt.Printf("Id %s not found. Aborting operation, nothing will be deleted\n", arg)
					return
				}
			}
		}

		// delete each of deletion groups passed as program arguments
		{
			for _, arg := range args {
				deletionPath := path.Join(trashPath, delGroupMap[arg])
				os.RemoveAll(deletionPath)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(flushCmd)
	flushCmd.Flags().Bool(flush_all_flag, false, "Permanently delete all files from trash")
}
