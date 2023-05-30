//go:build functional

package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/startupbuilder/startupbuilder/pkg/config"
)

type otherPtrTestConfig struct {
	Config *testConfig `mapstructure:"CONFIG"`
}

func (c *otherPtrTestConfig) Defaults() {
	c.Config = &testConfig{}
}

type otherTestConfig struct {
	Config testConfig `mapstructure:"config"`
}

func (c *otherTestConfig) Defaults() {
}

type testConfig struct {
	Demo      int    `mapstructure:"demo"`
	Multiname string `mapstructure:"multi.name"`
}

func (c *testConfig) Defaults() {}

func Test_OSEnvOneLevel(t *testing.T) {
	// Scenario one level struct
	cfg := &testConfig{}

	os.Setenv("EXTRA_DEMO", "1")
	os.Setenv("EXTRA_MULTI_NAME", "yesy")
	err := config.ReadConfig(cfg, config.WithPrefix("EXTRA"))
	require.NoError(t, err)

	require.Equal(t, 1, cfg.Demo)
	require.Equal(t, "yesy", cfg.Multiname)
}

func Test_OSEnvPtrTwoLevel(t *testing.T) {
	// Scenario two level struct
	otherCfg := otherPtrTestConfig{}

	os.Setenv("EXTPTR_CONFIG_DEMO", "2")
	os.Setenv("EXTPTR_CONFIG_MULTI_NAME", "22")
	err := config.ReadConfig(&otherCfg, config.WithPrefix("EXTPTR"))
	require.NoError(t, err)

	require.Equal(t, 2, otherCfg.Config.Demo)
	require.Equal(t, "22", otherCfg.Config.Multiname)
}

func Test_OSEnvTwoLevel(t *testing.T) {
	// Scenario two level struct
	otherCfg := otherTestConfig{}

	os.Setenv("EXT_CONFIG_DEMO", "2")
	os.Setenv("EXT_CONFIG_MULTI_NAME", "22")
	err := config.ReadConfig(&otherCfg, config.WithPrefix("EXT"))
	require.NoError(t, err)

	require.Equal(t, 2, otherCfg.Config.Demo)
	require.Equal(t, "22", otherCfg.Config.Multiname)
}

func Test_FileOneLevel(t *testing.T) {
	// Scenario one level struct
	cfg := &testConfig{}

	err := config.ReadConfig(cfg, config.WithFilePath(".env.test"))
	require.NoError(t, err)

	require.Equal(t, 12, cfg.Demo)
	require.Equal(t, "file", cfg.Multiname)
}

func Test_FilePtrTwoLevel(t *testing.T) {
	// Scenario two level struct
	otherCfg := otherPtrTestConfig{}

	err := config.ReadConfig(&otherCfg, config.WithFilePath(".env.test"))
	require.NoError(t, err)

	require.Equal(t, 123, otherCfg.Config.Demo)
	require.Equal(t, "file2", otherCfg.Config.Multiname)
}

func Test_FileTwoLevel(t *testing.T) {
	// Scenario two level struct
	otherCfg := otherTestConfig{}

	err := config.ReadConfig(&otherCfg, config.WithFilePath(".env.test"))
	require.NoError(t, err)

	require.Equal(t, 123, otherCfg.Config.Demo)
	require.Equal(t, "file2", otherCfg.Config.Multiname)
}

func Test_FileTwoLevelWithOverride(t *testing.T) {
	// Scenario two level struct
	otherCfg := otherTestConfig{}

	os.Setenv("OVER_CONFIG_DEMO", "1234")
	err := config.ReadConfig(&otherCfg, config.WithFilePath(".env.test"), config.WithPrefix("OVER"))
	require.NoError(t, err)

	require.Equal(t, 1234, otherCfg.Config.Demo)
	require.Equal(t, "file2", otherCfg.Config.Multiname)
}

func Test_OneLevelWithDefaults(t *testing.T) {
	// Should not override the default if no value set
	cfg := &testConfig{
		Demo:      2,
		Multiname: "123",
	}

	err := config.ReadConfig(cfg, config.WithPrefix("DEFAULTS"))
	require.NoError(t, err)

	require.Equal(t, 2, cfg.Demo)
	require.Equal(t, "123", cfg.Multiname)
}

func Test_OsEnvOneLevelWithDefaults(t *testing.T) {
	// Scenario one level struct mix between envs
	cfg := &testConfig{
		Demo:      2,
		Multiname: "123",
	}

	os.Setenv("OSDEFAULTS_DEMO", "1")
	os.Setenv("OSDEFAULTS_MULTI_NAME", "yesy")
	err := config.ReadConfig(cfg, config.WithPrefix("OSDEFAULTS"))
	require.NoError(t, err)

	require.Equal(t, 1, cfg.Demo)
	require.Equal(t, "yesy", cfg.Multiname)
}
