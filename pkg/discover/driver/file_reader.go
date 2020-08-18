package driver

import (
	"fmt"
	"io/ioutil"
)

// SourceFileReader used to read config source from a file
// implement discocver.SourceReader
type SourceFileReader struct{}

// NewFileReader used to create new instance of os file reader
// using ioutil as implementation
func NewFileReader() *SourceFileReader {
	return &SourceFileReader{}
}

func (r *SourceFileReader) Read(filepath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return b, nil
}
