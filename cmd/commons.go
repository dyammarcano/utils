package cmd

import (
	"os"
	"path/filepath"
)

func currentDirectory(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		if args[0] != "." {
			wd = filepath.Join(wd, args[0])
		}
	}

	return wd, nil
}
