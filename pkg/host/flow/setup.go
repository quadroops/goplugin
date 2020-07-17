package flow

import (
	"context"
	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/reactivex/rxgo/v2"
)

// Setup is main flows for host setup process
type Setup struct {
	Hostname string
	Config   *discover.PluginConfig
}

// NewSetup used to create new instance
func NewSetup(hostname string, config *discover.PluginConfig) *Setup {
	return &Setup{hostname, config}
}

// FilterByHostname used to filtering plugin config based on hostname
func (s *Setup) FilterByHostname(i interface{}) bool {
	config, ok := i.(*discover.PluginConfig)
	if !ok {
		return false
	}

	_, exist := config.Hosts[s.Hostname]
	return exist
}

// FilterByPlugin used to filtering plugin based on given plugin's name
// need to make sure if incoming plugin is exist
func (s *Setup) FilterByPlugin(i interface{}) bool {
	plugin, ok := i.(string)
	if !ok {
		return false
	}

	_, exist := s.Config.Plugins[plugin]
	return exist
}

// MapToPlugins used to map incoming data plugin config into list of plugin's name
func (s *Setup) MapToPlugins(_ context.Context, i interface{}) (interface{}, error) {
	config, ok := i.(*discover.PluginConfig)
	if !ok {
		return nil, errs.ErrCastInterface
	}

	plugins, exist := config.Hosts[s.Hostname]
	if !exist { return nil, errs.ErrNoPlugins }

	return plugins.Plugins, nil
}

// FlatToNewObservable used to create new observable timeline using Defer
// to pending creation of observable until all observers subscribed
func (s *Setup) FlatToNewObservable(i rxgo.Item) rxgo.Observable {
	plugins, ok := i.V.([]string)
	if !ok {
		return rxgo.Just(errs.ErrCastInterface)()
	}

	if len(plugins) < 1 {
		return rxgo.Just(errs.ErrNoPlugins)()
	}

	var newSources []rxgo.Producer
	source := func(_ context.Context, next chan<- rxgo.Item) {
		for _, plugin := range plugins {
			next <- rxgo.Of(plugin)
		}
	}

	newSources = append(newSources, source)
	return rxgo.Defer(newSources)
}
