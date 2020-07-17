package driver

import (
	"fmt"
	"io/ioutil"

	"github.com/quadroops/goplugin/pkg/discover"
)

type reader struct {}

// NewFileReader used to create new instance of os file reader
// using ioutil as implementation
func NewFileReader() discover.FileReader {
	return &reader{}
}

func (r *reader) ReadFile(filepath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return b, nil
}