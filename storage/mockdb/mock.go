// Package mock provides a mock storage Shipper interface for tests.
package mockdb

import (
	"errors"
	"github.com/tmaiaroto/discfg/config"
)

// MockCfg is just a map of mock records within a mock config.
var MockCfg = map[string]map[string]config.Item{
	"mockcfg": {
		"/": config.Item{
			Key:                    "/",
			Value:                  []byte("Mock configuration"),
			CfgVersion:             int64(4),
			CfgModifiedNanoseconds: int64(1464675792991825937),
		},
		"initial": config.Item{
			Key:     "initial",
			Value:   []byte("initial value for test"),
			Version: int64(1),
		},
		"initial_second": config.Item{
			Key:     "initial_second",
			Value:   []byte("a second initial value for test"),
			Version: int64(3),
		},
		"json_value": config.Item{
			Key:     "initial_second",
			Value:   []byte(`{"json": "string", "num": 4}`),
			Version: int64(3),
		},
		"encoded": config.Item{
			Key:     "encoded",
			Value:   []byte(`eyJ1cGRhdGVkIjogImZyaWRheSJ9`),
			Version: int64(1),
		},
	},
}

// Version int64  `json:"version,omitempty"`
// Key     string `json:"key,omitempty"`
// Value interface{} `json:"value,omitempty"`
// TTL              int64     `json:"ttl,omitempty"`
// Expiration       time.Time `json:"-"`
// OutputExpiration string    `json:"expiration,omitempty"`
// CfgVersion             int64 `json:"-"`
// CfgModifiedNanoseconds int64 `json:"-"`

// MockShipper struct implements the Shipper interface for testing purposes.
type MockShipper struct {
}

// Name returns the name for the interface
func (m MockShipper) Name(opts config.Options) string {
	return "Mock Storage Engine"
}

// Options returns various settings and options for the shipper
func (m MockShipper) Options(opts config.Options) map[string]interface{} {
	return map[string]interface{}{"example": "option"}
}

// CreateConfig creates a config
func (m MockShipper) CreateConfig(opts config.Options, settings map[string]interface{}) (interface{}, error) {
	var err error
	return "", err
}

// DeleteConfig deletes a config
func (m MockShipper) DeleteConfig(opts config.Options) (interface{}, error) {
	var err error
	return "", err
}

// UpdateConfig updates a config
func (m MockShipper) UpdateConfig(opts config.Options, settings map[string]interface{}) (interface{}, error) {
	var err error
	return "", err
}

// ConfigState returns the state of the config
func (m MockShipper) ConfigState(opts config.Options) (string, error) {
	var err error
	return "ACTIVE", err
}

// Update a Item (record)
func (m MockShipper) Update(opts config.Options) (config.Item, error) {
	var err error
	if val, ok := MockCfg[opts.CfgName][opts.Key]; ok {
		val.Version++
	} else {
		MockCfg[opts.CfgName][opts.Key] = config.Item{
			Key:     opts.Key,
			Value:   opts.Value,
			Version: int64(1),
		}
	}
	return MockCfg[opts.CfgName][opts.Key], err
}

// Get a Item (record)
func (m MockShipper) Get(opts config.Options) (config.Item, error) {
	var err error
	return MockCfg[opts.CfgName][opts.Key], err
}

// Delete a Item (record)
func (m MockShipper) Delete(opts config.Options) (config.Item, error) {
	var err error
	defer delete(MockCfg[opts.CfgName], opts.Key)
	return MockCfg[opts.CfgName][opts.Key], err
}

// UpdateConfigVersion updates the incremental counter/state of a configuration and should be called on each change
func (m MockShipper) UpdateConfigVersion(opts config.Options) error {
	var err error
	if opts.CfgName != "" {
		n := MockCfg[opts.CfgName]["/"]
		n.CfgVersion++
		MockCfg[opts.CfgName]["/"] = n
	} else {
		err = errors.New("Interface Error: No config name passed.")
	}
	return err
}
