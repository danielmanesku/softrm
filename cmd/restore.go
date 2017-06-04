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
	"path/filepath"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore ID [ID...]",
	Short: "Restore removed items",
	Long: `Restore given deletion groups back to their original location.
Use listing command to find deletion group IDs of files
from trash.`,
	Run: func(cmd *cobra.Command, args []string) {

		// validate inputs
		{
			if len(args) < 1 {
				cmd.Help()
				return
			}
		}

		ensureTrashDirExists(false)
		trashPath := getTrashPath()

		delGroupMap := getDeletionGroupIds(trashPath)
		checkValidityOfIdArgs(args, delGroupMap)

		// restore each of deletion groups passed as program arguments
		{
			for _, arg := range args {
				argStartingPath, _ := filepath.Abs(path.Join(trashPath, delGroupMap[arg]))
				curPath := argStartingPath

				// collect all path segments
				var restorePathSegments []string
				for {
					if files, err := ioutil.ReadDir(curPath); err != nil {
						dirReadError(curPath, err)
					} else {
						if len(files) == 1 && files[0].IsDir() { // walk up through directories
							// as long as there is only one child directory
							segmentName := files[0].Name()
							restorePathSegments = append(restorePathSegments, segmentName)
							curPath = path.Join(curPath, segmentName)
						} else {
							break // all segments collected
						}
					}
				}

				restoreRoot := "/" + path.Join(restorePathSegments...)
				ensureDirExists(restoreRoot, true)

				// move each file/dir to restore destination root
				itemsToRestore, err := ioutil.ReadDir(curPath)
				if err != nil {
					dirReadError(curPath, err)
				}
				for _, item := range itemsToRestore {
					if err := os.Rename(path.Join(curPath, item.Name()), path.Join(restoreRoot, item.Name())); err != nil {
						fmt.Println(err.Error())
						fmt.Println("Note: some files might have already been moved. Check destination directory")
						abortAndExit()
					}
				}

				// remove root directory for given arg (deletion group id)
				os.RemoveAll(argStartingPath)

				fmt.Printf("Restoration of %s completed successfully\n", arg)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)
}
