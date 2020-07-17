package driver

import (
	"fmt"
	"io"
	"os"
	"crypto/md5"
	"github.com/quadroops/goplugin/pkg/host"
) 

type driverMd5Checker struct {}

// NewMd5Check used to create an instance of driver md5 checker
func NewMd5Check() host.MD5Checker {
	return &driverMd5Checker{}
}

func (drv *driverMd5Checker) Parse(file string) (string, error) {
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
