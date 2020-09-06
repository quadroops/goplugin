package goplugin

import (
	"github.com/quadroops/goplugin/internal/factory"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"
)

// GoPlugin .
type GoPlugin struct {
	hostName        string
	configChecker   *discover.ConfigChecker
	configParser    *discover.ConfigParser
	processInstance *process.Instance
	identityChecker host.IdentityChecker
}

// Option used to customize default objects
type Option func(*GoPlugin)

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
