package print

import (
	"fmt"
	"os"
	"path/filepath"
)

func printFileSystemTree(path, indent string) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, dir := range dirs {
		isLast := i == len(dirs)-1
		printDir(dir.Name(), indent, isLast)

		if dir.IsDir() {
			lineStyle := "│   "
			if isLast {
				lineStyle = "    "
			}
			subPath := filepath.Join(path, dir.Name())
			if err := printFileSystemTree(subPath, indent+lineStyle); err != nil {
				return err
			}
		}
	}

	return nil
}

func printDir(dirName, indent string, isLast bool) {
	prefix := "├── "
	if isLast {
		prefix = "└── "
	}
	fmt.Printf("%s%s%s\n", indent, prefix, dirName)
}

func PrintDirectoryTree(rootPath string) error {
	fmt.Println(rootPath)
	return printFileSystemTree(rootPath, "")
}
