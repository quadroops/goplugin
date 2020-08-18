package flow

import "os"

// IdentityCheckerProxy used as proxy interface to solve
// cyclic dependency relate with MD5Checker
type IdentityCheckerProxy interface {
	Parse(file string) (string, error)
}

// RegistryProxy used as proxy to host.Registry
type RegistryProxy struct {
	ExecPath     string
	ExecArgs     []string
	ExecFile     string
	ExecTime     int
	MD5Sum       string
	ProtocolType string
}

// Plugin as main observable item
type Plugin struct {
	Name     string
	Registry RegistryProxy
}

// Install used as main flows for install process
type Install struct {
	IDChecker IdentityCheckerProxy
}

// NewInstall used to create new instance
func NewInstall(checker IdentityCheckerProxy) *Install {
	return &Install{checker}
}

// FilterByExecFile used to filtering item by checking plugin's file
// existence on os
func (i *Install) FilterByExecFile(v interface{}) bool {
	plugin, ok := v.(Plugin)
	if !ok {
		return false
	}

	_, err := os.Stat(plugin.Registry.ExecFile)
	return err == nil
}

// FilterByMD5 used to filtering item by checking plugin's md5sum
func (i *Install) FilterByMD5(v interface{}) bool {
	plugin, ok := v.(Plugin)
	if !ok {
		return false
	}

	md5Str, err := i.IDChecker.Parse(plugin.Registry.ExecFile)
	if err != nil {
		return false
	}

	return md5Str == plugin.Registry.MD5Sum
}
