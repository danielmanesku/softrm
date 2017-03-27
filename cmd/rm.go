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
	"strings"
	"time"

	"github.com/danielmanesku/softrm/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Move file(s) to trash",
	Long: `Move files(s) to trash.
It requires one or more additional parameters that are file
or directory names in current working directory. Absolute
paths are supported as well.

For example:
softrm rm tempfile.txt`,
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
		ensureTrashDirExists()

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
			// create directory structure for new file/dir, since os.Rename can't do it
			err := os.MkdirAll(path.Join(delInstancePath, filepath.Dir(argPath)), 0700)
			if nil != err {
				fmt.Println(err.Error())
				abortAndExit()
			}

			// now move the file/dir
			if err := os.Rename(argPath, path.Join(delInstancePath, argPath)); err != nil {
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

// Check if trash directory exists and has adequate permissions. Create one
// if it doesn't.
func ensureTrashDirExists() {
	trashPath := getTrashPath()
	fi, err := os.Stat(trashPath)
	if err == nil {
		if fi.IsDir() {
			if 0 != strings.Compare(fi.Mode().Perm().String()[1:4], "rwx") {
				fmt.Printf("Directory %s does not have 'rwx' permissions for current user.\n", trashPath)
				abortAndExit()
			}
			return // all good
		} else {
			fmt.Printf("Node %s exists, but it's not a directory.\n", trashPath)
			abortAndExit()
		}
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(trashPath, 0700)
		if nil == err {
			fmt.Printf("Directory %s created.\n", trashPath)
			return
		} else {
			fmt.Println(err.Error())
			abortAndExit()
		}
	}
	fmt.Println("error:", err)
	abortAndExit()
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

// Trash path is currently pointing to config value.
// Later it might be more trash paths, one for each partition, to avoid
// unnecessary data transfers if source and trash path are not on
// the same partition.
func getTrashPath() string {
	return configTrashPath()
}

func configTrashPath() string {
	return os.ExpandEnv(viper.GetString("trashdir"))
}

func abortAndExit() {
	//TODO maybe replace with log.Fatal
	fmt.Println("Aborting.")
	os.Exit(1)
}
