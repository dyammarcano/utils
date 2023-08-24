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

func NewMockFile(content []byte) *MockFile {
	return &MockFile{
		Content: content,
	}
}

func (mf *MockFile) Read(p []byte) (n int, err error) {
	if mf.Pos >= len(mf.Content) {
		return 0, io.EOF
	}
	n = copy(p, mf.Content[mf.Pos:])
	mf.Pos += n
	return n, nil
}

func (mf *MockFile) Write(p []byte) (n int, err error) {
	if mf.Pos >= len(mf.Content) {
		mf.Content = append(mf.Content, p...)
	} else {
		n = copy(mf.Content[mf.Pos:], p)
		mf.Pos += n
	}
	return len(p), nil
}

func (mf *MockFile) Close() error {
	return nil
}

func CreateTmpDir(t *testing.T) (string, error) {
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
		return "", err
	}
	return dir, nil
}

func CleanUpTmpDir(t *testing.T, path string) {
	<-time.After(5 * time.Second)
	if err := os.RemoveAll(path); err != nil {
		t.Fatal("Failed to cleanup temp directory:", err)
	}
}

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

func NewLoremTestFile(t *testing.T, phrases int) (*Lorem, error) {
	lorem := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.")
	var data []byte

	if phrases <= 0 {
		phrases = 1
	}

	for i := 0; i < phrases; i++ {
		data = append(data, lorem...)
	}

	dir, err := CreateTmpDir(t)
	if err != nil {
		return nil, err
	}

	testfile := filepath.Join(dir, "data.bin")

	fw, err := os.OpenFile(testfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("Could not open test file for writing:", err)
		return nil, err
	}

	if _, err = fw.Write(data); err != nil {
		t.Fatal("Failed to write to test file:", err)
		return nil, err
	}

	return &Lorem{
		LoremFilePath: testfile,
		d:             dir,
		t:             t,
		File:          fw,
	}, nil
}

func (l *Lorem) CleanUp() {
	if err := l.File.Close(); err != nil {
		l.t.Fatal("Failed to close test file:", err)
	}

	<-time.After(100 * time.Millisecond)

	if err := os.RemoveAll(l.d); err != nil {
		l.t.Fatal("Failed to cleanup temp directory:", err)
	}
}
