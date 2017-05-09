package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Check if trash directory exists and has adequate permissions. Create one
// if it doesn't.
func ensureTrashDirExists(verbose bool) {
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
		createDir(trashPath)
		if verbose {
			fmt.Printf("Directory %s created.\n", trashPath)
		}
		return
	}
	fmt.Println("error:", err)
	abortAndExit()
}

func createDir(path string) {
	err := os.MkdirAll(path, 0700)
	if nil != err {
		fmt.Println(err.Error())
		abortAndExit()
	}
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
