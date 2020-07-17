package discover

import (
	"fmt"

	"github.com/quadroops/goplugin/pkg/errs"
)

// ConfigParser used to parse configuration values from given filepath 
type ConfigParser struct {
	parser Parser
	reader FileReader
}

// NewConfigParser used to create new instance of ConfigParser
func NewConfigParser(parser Parser, reader FileReader) *ConfigParser {
	return &ConfigParser{parser: parser, reader: reader}
}

// Load used to read given config's file, get the content bytes and parse the data 
func (cp *ConfigParser) Load(confpath string) (*PluginConfig, error) {
	content, err := cp.reader.ReadFile(confpath) 
	if err != nil {
		return nil, fmt.Errorf("Confpath: %q %w", confpath, errs.ErrReadConfigFile) 
	}

	conf, err := cp.parser.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("%w", errs.ErrParseConfig)
	}

	return conf, nil
}