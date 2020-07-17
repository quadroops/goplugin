package flow

import "os"

// MD5CheckerProxy used as proxy interface to solve
// cyclic dependency relate with MD5Checker
type MD5CheckerProxy interface {
	Parse(file string) (string, error)
}

// RegistryProxy used as proxy to host.Registry
type RegistryProxy struct {
	ExecPath string
	ExecArgs []string
	ExecFile string
	ExecTime int
	MD5Sum   string
	RPCType  string
	RPCPort  string
}

// Plugin as main observable item
type Plugin struct {
	Name     string
	Registry RegistryProxy
}

// Install used as main flows for install process
type Install struct {
	MD5Checker MD5CheckerProxy
}

// NewInstall used to create new instance 
func NewInstall(md5Checker MD5CheckerProxy) *Install {
	return &Install{md5Checker}
}

// FilterByExecFile used to filtering item by checking plugin's file
// existence on os
func (i *Install) FilterByExecFile(v interface{}) bool {
	plugin, ok := v.(Plugin)
	if !ok { return false }

	_, err := os.Stat(plugin.Registry.ExecFile)
	return err == nil
}

// FilterByMD5 used to filtering item by checking plugin's md5sum
func (i *Install) FilterByMD5(v interface{}) bool {
	plugin, ok := v.(Plugin)
	if !ok { return false }

	md5Str, err := i.MD5Checker.Parse(plugin.Registry.ExecFile)
	if err != nil { return false }

	return md5Str == plugin.Registry.MD5Sum
}