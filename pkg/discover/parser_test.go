package discover_test

import (
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/discover/mocks"
	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func _generate_config() *discover.PluginConfig {
	meta := discover.PluginMeta{
		Author:       "test-author",
		Contributors: []string{"test1"},
		Version:      "1.0.0",
	}

	settings := discover.PluginSettings{
		Debug: true,
	}

	pluginMap := make(map[string]discover.PluginInfo)
	pluginMap["test"] = discover.PluginInfo{
		Author:       "test-author-plugin",
		Exec:         "test-exec",
		ExecTime:     10,
		MD5:          "2323233",
		ProtocolType: "nano",
	}

	hostMap := make(map[string]discover.PluginHost)
	hostMap["test-host"] = discover.PluginHost{
		Plugins: []string{"test"},
	}

	return &discover.PluginConfig{
		Meta:     meta,
		Settings: settings,
		Hosts:    hostMap,
		Plugins:  pluginMap,
	}
}

func TestLoadConfigSuccess(t *testing.T) {
	confPath := "config.toml"

	mockTomlParser := new(mocks.Parser)
	mockSourceReader := new(mocks.SourceReader)

	mockSourceReader.On("Read", confPath).Return([]byte("test"), nil)
	mockTomlParser.On("Parse", []byte("test")).Return(_generate_config(), nil)

	configParser := discover.NewConfigParser(mockTomlParser, mockSourceReader)
	config, err := configParser.Load(confPath)
	assert.NoError(t, err)
	assert.Equal(t, config.Meta.Version, "1.0.0")
	assert.True(t, mockSourceReader.AssertCalled(t, "Read", confPath))
	assert.True(t, mockTomlParser.AssertCalled(t, "Parse", []byte("test")))
}

func TestLoadConfigErrorFileReader(t *testing.T) {
	confPath := "config.toml"

	mockTomlParser := new(mocks.Parser)
	mockSourceReader := new(mocks.SourceReader)

	mockSourceReader.On("Read", confPath).Return(nil, errors.New("error reader"))
	configParser := discover.NewConfigParser(mockTomlParser, mockSourceReader)
	_, err := configParser.Load(confPath)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrReadConfigFile))

	assert.True(t, mockSourceReader.AssertCalled(t, "Read", confPath))
	assert.True(t, mockTomlParser.AssertNotCalled(t, "Parse", []byte("test")))
}

func TestLoadConfigErrorParser(t *testing.T) {
	confPath := "config.toml"

	mockTomlParser := new(mocks.Parser)
	mockSourceReader := new(mocks.SourceReader)

	mockSourceReader.On("Read", confPath).Return([]byte("test"), nil)
	mockTomlParser.On("Parse", []byte("test")).Return(nil, errors.New("error parser"))
	configParser := discover.NewConfigParser(mockTomlParser, mockSourceReader)
	_, err := configParser.Load(confPath)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrParseConfig))

	assert.True(t, mockSourceReader.AssertCalled(t, "Read", confPath))
	assert.True(t, mockTomlParser.AssertCalled(t, "Parse", []byte("test")))
}
