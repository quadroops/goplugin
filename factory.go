package goplugin

import (
	"github.com/quadroops/goplugin/internal/factory"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"
)

// WithCustomConfigChecker used to customize config checker, the parameter
// must be implement discover.Checker
func WithCustomConfigChecker(adapters ...discover.Checker) Option {
	return func(gp *GoPlugin) {
		gp.configChecker = discover.NewConfigChecker(adapters...)
	}
}

// WithCustomConfigParser used to customize config reader & parser, given adapters
// must be implement discover.Parser and discover.SourceReader interfaces
func WithCustomConfigParser(parserAdapter discover.Parser, sourceAdapter discover.SourceReader) Option {
	return func(gp *GoPlugin) {
		gp.configParser = discover.NewConfigParser(parserAdapter, sourceAdapter)
	}
}

// WithCustomIdentityChecker used to customize plugin's identity checker, given adapter
// must implement host.IdentityChecker
func WithCustomIdentityChecker(adapter host.IdentityChecker) Option {
	return func(gp *GoPlugin) {
		gp.identityChecker = adapter
	}
}

// WithCustomProcess used to customize process runner and process data builder (including registry)
func WithCustomProcess(runner process.Runner, processes process.ProcessesBuilder) Option {
	return func(gp *GoPlugin) {
		gp.processInstance = process.New(runner, processes)
	}
}

// Map used to put a plugin and assign it with their spesific configurations
func Map(pluginName string, conf *PluginConf) PluginMapper {
	mapper := make(PluginMapper)
	mapper[pluginName] = conf
	return mapper
}

// New used to create new instance of GoPlugin, you can provide customization
// here by giving goplugin.Option
func New(hostName string, opts ...Option) *GoPlugin {
	gp := &GoPlugin{
		hostName:        hostName,
		configChecker:   factory.DefaultConfigChecker(),
		configParser:    factory.DefaultConfigParser(),
		processInstance: factory.DefaultProcessInstance(),
		identityChecker: factory.DefaultHostIdentityChecker(),
	}

	for _, option := range opts {
		option(gp)
	}

	return gp
}

// SetupPlugin used to store host plugin's configurations
func (g *GoPlugin) SetupPlugin(plugins ...PluginMapper) *GoPlugin {
	mapper := make(PluginMapper)

	// we're now need to merge all given plugin's config mapper
	for _, plugin := range plugins {
		for k, v := range plugin {
			mapper[k] = v
		}
	}

	g.hostPlugins = mapper
	return g
}

// GetPluginConf used to get plugin's config
func (g *GoPlugin) GetPluginConf(pluginName string) (*PluginConf, error) {
	if conf, exist := g.hostPlugins[pluginName]; exist {
		return conf, nil
	}

	return nil, errs.ErrPluginNotFound
}

// Build used to compile current host instance into consumed state of host.Builder
func (g *GoPlugin) Build() (*host.Builder, error) {
	configAddr, err := g.configChecker.Explore()
	if err != nil {
		return nil, err
	}

	config, err := g.configParser.Load(configAddr)
	if err != nil {
		return nil, err
	}

	h := host.New(g.hostName, config, g.identityChecker)
	return h, nil
}

// GetProcessInstance .
func (g *GoPlugin) GetProcessInstance() *process.Instance {
	return g.processInstance
}
