package errs 

import "errors"

var (
	// ErrCastInterface used when failed to cast some interface value
	ErrCastInterface = errors.New("Cannot cast interface")

	// ErrReadConfigFile used when failed to read given config filepath
	ErrReadConfigFile = errors.New("Cannot read given config path")

	// ErrParseConfig used when cannot parse toml's contents
	ErrParseConfig = errors.New("Failed to parse configurations")

	// ErrEmptyPlugins used when host trying install with empty Plugins
	ErrEmptyPlugins = errors.New("Cannot install empty plugins")

	// ErrNoPlugins used when host failed to validate their available plugins
	ErrNoPlugins = errors.New("No plugins available")

	// ErrConfigNotFound used when failed to explore given file
	ErrConfigNotFound = errors.New("Config file not found")

	// ErrPluginNotFound used when cannot found requested plugin
	ErrPluginNotFound = errors.New("Plugin not found")

	// ErrPluginCannotStart used when failed to run a subprocess for given plugin
	ErrPluginCannotStart = errors.New("Plugin cannot start")

	// ErrPluginCannotBeKilled used when failing to kill the plugin
	ErrPluginCannotBeKilled = errors.New("Plugin cannot be killed")

	// ErrPluginStarted used when host try to run a plugin twice
	ErrPluginStarted = errors.New("Plugin has been started")

	// ErrEmptyProcesses used when there are no processes attached
	ErrEmptyProcesses = errors.New("No processes available")

	// ErrPluginPing used when cannot connect to plugin
	ErrPluginPing = errors.New("Plugin cannot be called")

	// ErrPluginCall used when found any errors when calling a plugin
	ErrPluginCall = errors.New("Plugin communication error")

	// ErrPluginExec used when found any errors when calling an exec command to plugin
	ErrPluginExec = errors.New("Plugin cannot exec")
	
	// ErrProtocolUnknown used when plugin define unsuppported protocol
	ErrProtocolUnknown = errors.New("Illegal protocol")

	// ErrProtocolGRPCConnection used for an error grpc connection
	ErrProtocolGRPCConnection = errors.New("Error grpc connection")

	// ErrProtocolGRPCResponse used when an error on response
	ErrProtocolGRPCResponse = errors.New("Error grpc response")
) 