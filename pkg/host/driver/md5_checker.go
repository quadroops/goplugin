package driver

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// Md5Checker used to parse md5 value from a given file path
// implement host.IdentityChecker
type Md5Checker struct{}

// NewMd5Check used to create an instance of driver md5 checker
func NewMd5Check() *Md5Checker {
	return &Md5Checker{}
}

// Parse get the md5 value from a file
func (drv *Md5Checker) Parse(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
