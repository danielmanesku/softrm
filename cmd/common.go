package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type DeletionGroup struct {
	id           string
	startingDir  string
	deletionTime time.Time
}

// Check if given directory exists and has adequate permissions. Create one
// if it doesn't.
func ensureDirExists(dir string, verbose bool) {
	fi, err := os.Stat(dir)
	if err == nil {
		if fi.IsDir() {
			if 0 != strings.Compare(fi.Mode().Perm().String()[1:4], "rwx") {
				fmt.Printf("Directory %s does not have 'rwx' permissions for current user.\n", dir)
				abortAndExit()
			}
			return // all good
		} else {
			fmt.Printf("Node %s exists, but it's not a directory.\n", dir)
			abortAndExit()
		}
	}
	if os.IsNotExist(err) {
		createDir(dir)
		if verbose {
			fmt.Printf("Directory %s created.\n", dir)
		}
		return
	}
	fmt.Println("error:", err)
	abortAndExit()
}

func ensureTrashDirExists(verbose bool) {
	trashPath := getTrashPath()
	ensureDirExists(trashPath, verbose)
}

func createDir(path string) {
	err := os.MkdirAll(path, 0700)
	if nil != err {
		fmt.Println(err.Error())
		abortAndExit()
	}
}

// getSelectedDeletionGroups will try to match all given (short) ids against
// existing ids from trash path. It will exit if any of those is not found.
// It is possible to use as id only a prefix of deletion group. If there is
// ambiguity (more than one match), command will print an error and exit.
func getSelectedDeletionGroups(ids []string, trashPath string) []DeletionGroup {
	delGroups := make([]DeletionGroup, len(ids))
	found := 0

	trashDirs, _ := ioutil.ReadDir(trashPath)
	for _, dir := range trashDirs {
		dirName := dir.Name()
		// take last part from dir name, after last dash
		idPart := dirName[strings.LastIndex(dirName, "-")+1:]
		for idx, curId := range ids {
			if strings.HasPrefix(idPart, curId) { // current directory from
				// trashPath starts with one of ids we are evaluating

				if delGroups[idx].id != "" {
					fmt.Printf("ID \"%s\" has more than one match. Please specify more characters.\n", curId)
					abortAndExit()
				}

				delGr := DeletionGroup{
					id:          idPart,
					startingDir: dirName,
				}
				delGroups[idx] = delGr
				found += 1
			}
		}
	}

	// check if all given ids were found
	if found != len(ids) {
		for idx, curId := range ids {
			if delGroups[idx].id == "" {
				fmt.Printf("Id %s not found.\n", curId)
				abortAndExit()
			}
		}
	}

	return delGroups
}

func dirReadError(dir string, err error) {
	fmt.Println("Error occurred: could not read contents of directory:", dir)
	fmt.Println("Original error is:", err)
	abortAndExit()
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
	fmt.Println("Aborting.")
	os.Exit(1)
}
