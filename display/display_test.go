package display

import (
	"github.com/dyammarcano/utils/mocks"
	"os"
	"path/filepath"
	"testing"
)

func TempDirectoryWithFiles(t *testing.T) (dirPath string, cleanup func()) {
	t.Helper()

	dir, err := mocks.CreateTempDir(t)
	if err != nil {
		t.Errorf("Failed to create tmp dir: %v", err)
	}

	// make second level directory
	dir2 := filepath.Join(dir, "dir2")
	if err := os.Mkdir(dir2, 0755); err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	// make third level directory
	dir3 := filepath.Join(dir2, "dir3")
	if err := os.Mkdir(dir3, 0755); err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	files := []string{
		"file1.txt",
		"file2.txt",
	}

	for _, file := range files {
		f, err := os.Create(filepath.Join(dir, file))
		if err != nil {
			t.Fatalf("failed to create temp file: %v, %v", file, err)
		}

		if _, err := f.Write([]byte("This is a test file")); err != nil {
			if err := f.Close(); err != nil {
				return "", nil
			}
			t.Fatalf("failed to write into the file: %v", err)
		}
		if err := f.Close(); err != nil {
			return "", nil
		}
	}

	return dir, func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatalf("failed to cleanup temp directory: %v", err)
		}
	}
}

func TestPrintDirectoryTree(t *testing.T) {
	dir, cleanup := TempDirectoryWithFiles(t)
	defer cleanup()

	if err := DisplayDirectoryTree(dir); err != nil {
		t.Fatal(err)
	}
}
