// Copyright © 2017 Daniel Manesku daniel.manesku@gmail.com
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
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/danielmanesku/softrm/util"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm FILE [FILE...]",
	Short: "Move file(s) to trash",
	Long: `Move files(s) to trash.
It requires one or more additional parameters that are file
or directory names in current working directory. Absolute
paths are supported as well.`,
	Example: "softrm rm tempfile.txt",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		// check if all arguments are valid path
		for _, argPath := range args {
			if _, err := os.Stat(argPath); os.IsNotExist(err) {
				fmt.Printf("Argument %s is not a valid path.\n", argPath)
				abortAndExit()
			}
		}
		ensureTrashDirExists(true)

		// create deletion instance directory
		var delInstancePath string
		{
			delInstanceName := genDeletionInstanceName()
			delInstancePath = path.Join(getTrashPath(), delInstanceName)
			err := os.MkdirAll(delInstancePath, 0700)
			if nil != err {
				fmt.Println(err.Error())
				abortAndExit()
			}
		}

		// move all files to deletion instance directory
		for _, argPath := range args {
			argAbsPath, _ := filepath.Abs(argPath)
			argDirPath := filepath.Dir(argAbsPath)

			// create directory structure for new file/dir, since os.Rename can't do it
			err := os.MkdirAll(path.Join(delInstancePath, argDirPath), 0700)
			if nil != err {
				fmt.Println(err.Error())
				abortAndExit()
			}

			// now move the file/dir
			if err := os.Rename(argPath, path.Join(delInstancePath, argAbsPath)); err != nil {
				fmt.Println(err.Error())
				abortAndExit()
			}
		}

		fmt.Println("All files were moved to directory", delInstancePath)
		fmt.Println("rm operation successfully completed.")
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}

// Deletion instance name (directory name) consists of current
// time (date and time with seconds precision) and unique id
// generated from nanoseconds time and encoded with 36 base (which
// is reversed to get different prefixes)
func genDeletionInstanceName() string {
	nanoTime := time.Now().UTC().UnixNano()
	id := util.Reverse(strconv.FormatInt(nanoTime, 36))
	return time.Now().Format("2006-01-02T15-04-05") + "-" + id
}
