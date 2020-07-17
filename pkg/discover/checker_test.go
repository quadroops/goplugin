package discover_test

import (
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/discover/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSuccessLoadFromOS(t *testing.T) {
	mockReturnFile := "config_mock_os.toml"
	mockOSChecker := new(mocks.Checker)
	mockDefaultChecker := new(mocks.Checker)

	mockOSChecker.On("Check").Once().Return(mockReturnFile)
	mockDefaultChecker.On("Check").Return("")

	checker := discover.NewConfigChecker(mockOSChecker, mockDefaultChecker)
	conf, err := checker.Explore()
	assert.NoError(t, err)
	assert.Equal(t, conf, mockReturnFile)
	mockOSChecker.AssertCalled(t, "Check")
	mockDefaultChecker.AssertCalled(t, "Check")
}

func TestSuccessLoadFromDefault(t *testing.T) {
	mockReturnFile := "config_mock_default.toml"
	mockOSChecker := new(mocks.Checker)
	mockDefaultChecker := new(mocks.Checker)

	mockOSChecker.On("Check").Once().Return("")
	mockDefaultChecker.On("Check").Once().Return(mockReturnFile)

	checker := discover.NewConfigChecker(mockOSChecker, mockDefaultChecker)
	conf, err := checker.Explore()
	assert.NoError(t, err)
	assert.Equal(t, conf, mockReturnFile)
	mockOSChecker.AssertCalled(t, "Check")
	mockDefaultChecker.AssertCalled(t, "Check")
}

func TestErrorNotFound(t *testing.T) {
	mockOSChecker := new(mocks.Checker)
	mockDefaultChecker := new(mocks.Checker)

	mockOSChecker.On("Check").Once().Return("")
	mockDefaultChecker.On("Check").Once().Return("")

	checker := discover.NewConfigChecker(mockOSChecker, mockDefaultChecker)
	_, err := checker.Explore()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrConfigNotFound))
	mockOSChecker.AssertCalled(t, "Check")
	mockDefaultChecker.AssertCalled(t, "Check")
}