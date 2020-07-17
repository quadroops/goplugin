package driver

import (
	"fmt"

	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/pelletier/go-toml"
)

// TomlParser used as driver that implement Parse
type parser struct {}

// NewTomlParser create new instance of parser
func NewTomlParser() discover.Parser {
	return &parser{}
}

func (p *parser) Parse(content []byte) (*discover.PluginConfig, error) {
	var conf discover.PluginConfig
	tomlConf, err := toml.LoadBytes(content)
	if err != nil {
		return nil, fmt.Errorf("%v", err) 
	}

	err = tomlConf.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
