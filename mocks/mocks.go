package mocks

import "io"

type MockFile struct {
	Content []byte
	Pos     int
}

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
