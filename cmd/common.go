package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

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

// getDeletionGroupIds returns a map of all available deletion group ids and their directory names.
func getDeletionGroupIds(trashPath string) map[string]string {
	delGroupMap := make(map[string]string)
	files, _ := ioutil.ReadDir(trashPath)
	for _, f := range files {
		fileName := f.Name()
		// take last part from dir name, after last dash
		delGroupId := fileName[strings.LastIndex(fileName, "-")+1:]
		delGroupMap[delGroupId] = fileName
	}
	return delGroupMap
}

// checkValidityOfIdArgs will check if all given ids exist.
// In case any of given ids does not exist message will be printed and application will exit.
func checkValidityOfIdArgs(ids []string, delGroupMap map[string]string) {
	for _, id := range ids {
		_, found := delGroupMap[id]
		if !found {
			fmt.Printf("Id %s not found. Aborting operation, nothing will be deleted\n", id)
			abortAndExit()
		}
	}
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
