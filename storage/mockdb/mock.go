// Package mock provides a mock storage Shipper interface for tests.
package mockdb

import (
	"errors"
	"github.com/tmaiaroto/discfg/config"
)

// MockCfg is just a map of mock records within a mock config.
var MockCfg = map[string]map[string]config.Node{
	"mockcfg": {
		"/": config.Node{
			Key:                    "/",
			Value:                  []byte("Mock configuration"),
			CfgVersion:             int64(4),
			CfgModifiedNanoseconds: int64(1464675792991825937),
		},
		"initial": config.Node{
			Key:     "initial",
			Value:   []byte("initial value for test"),
			Version: int64(1),
		},
		"initial_second": config.Node{
			Key:     "initial_second",
			Value:   []byte("a second initial value for test"),
			Version: int64(3),
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
	return "", err
}

// Update a Node (record)
func (m MockShipper) Update(opts config.Options) (config.Node, error) {
	var err error
	if val, ok := MockCfg[opts.CfgName][opts.Key]; ok {
		val.Version++
	} else {
		MockCfg[opts.CfgName][opts.Key] = config.Node{
			Key:     opts.Key,
			Value:   opts.Value,
			Version: int64(1),
		}
	}
	return MockCfg[opts.CfgName][opts.Key], err
}

// Get a Node (record)
func (m MockShipper) Get(opts config.Options) (config.Node, error) {
	var err error
	return MockCfg[opts.CfgName][opts.Key], err
}

// Delete a Node (record)
func (m MockShipper) Delete(opts config.Options) (config.Node, error) {
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
		err = errors.New("No config name passed.")
	}
	return err
}
