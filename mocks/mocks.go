package mocks

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type (
	MockFile struct {
		Content []byte
		Pos     int
	}

	Lorem struct {
		d string
		t *testing.T
		*os.File
		LoremFilePath string
	}
)

// NewMockFile creates a new mock file.
func NewMockFile(content []byte) *MockFile {
	return &MockFile{
		Content: content,
		Pos:     0,
	}
}

// Read reads data from a mock file.
func (mf *MockFile) Read(p []byte) (n int, err error) {
	if mf.Pos >= len(mf.Content) {
		return 0, io.EOF
	}
	n = copy(p, mf.Content[mf.Pos:])
	mf.Pos += n
	return n, nil
}

// Write writes data to a mock file.
func (mf *MockFile) Write(p []byte) (n int, err error) {
	if mf.Pos >= len(mf.Content) {
		mf.Content = append(mf.Content, p...)
	} else {
		n = copy(mf.Content[mf.Pos:], p)
		mf.Pos += n
	}
	return len(p), nil
}

// CreateTempDir creates a temporary directory and returns the path to it.
func CreateTempDir(t *testing.T) (string, error) {
	dir, err := os.MkdirTemp("", "mock_")
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
		return "", err
	}
	return dir, nil
}

// CreateTempDirCleanUp creates a temporary directory and returns the path to it, an error and a cleanup function.
func CreateTempDirCleanUp(t *testing.T) (string, error, func()) {
	dir, err := CreateTempDir(t)
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
		return "", err, nil
	}

	return dir, nil, func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal("Failed to cleanup temp directory:", err)
		}
	}
}

// CleanUpTmpDir removes the temporary directory after 5 seconds.
func CleanUpTmpDir(t *testing.T, path string) {
	<-time.After(5 * time.Second)
	if err := os.RemoveAll(path); err != nil {
		t.Fatal("Failed to cleanup temp directory:", err)
	}
}

// WriteToTestFile writes data to a test file.
func WriteToTestFile(t *testing.T, testfile string, data []byte) error {
	fw, err := os.OpenFile(testfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("Could not open test file for writing:", err)
		return err
	}

	defer func(fw *os.File) {
		if err := fw.Close(); err != nil {
			t.Fatal("Failed to close test file:", err)
		}
	}(fw)

	if _, err = fw.Write(data); err != nil {
		t.Fatal("Failed to write to test file:", err)
		return err
	}

	return nil
}

// ReadFromTestFile reads data from a test file.
func ReadFromTestFile(t *testing.T, testfile string) ([]byte, error) {
	fo, err := os.Open(testfile)
	if err != nil {
		t.Fatalf("Could not open test file for reading: %v", err)
		return nil, err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			t.Fatalf("Failed to close test file: %v", err)
		}
	}()

	readData, err := io.ReadAll(fo)
	if err != nil {
		t.Fatalf("Failed to read from test file: %v", err)
		return nil, err
	}

	return readData, nil
}

// NewLoremTestFile creates a test file with lorem ipsum text.
func NewLoremTestFile(t *testing.T, phrases int) (*Lorem, error, func()) {
	lorem := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.")
	var data []byte

	if phrases <= 0 {
		phrases = 1
	}

	for i := 0; i < phrases; i++ {
		data = append(data, lorem...)
	}

	dir, err := CreateTempDir(t)
	if err != nil {
		return nil, err, nil
	}

	testfile := filepath.Join(dir, "data.bin")

	fw, err := os.OpenFile(testfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("Could not open test file for writing:", err)
		return nil, err, nil
	}

	if _, err = fw.Write(data); err != nil {
		t.Fatal("Failed to write to test file:", err)
		return nil, err, nil
	}

	return &Lorem{
			LoremFilePath: testfile,
			d:             dir,
			t:             t,
			File:          fw,
		}, nil, func() {
			if err := fw.Close(); err != nil {
				t.Fatal("Failed to close test file:", err)
			}

			<-time.After(100 * time.Millisecond)

			if err := os.RemoveAll(dir); err != nil {
				t.Fatal("Failed to cleanup temp directory:", err)
			}
		}
}
