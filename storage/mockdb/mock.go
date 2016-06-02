// Package mock provides a mock storage Shipper interface for tests.
package mockdb

import (
	"errors"
	"github.com/tmaiaroto/discfg/config"
)

// MockNodes is just a map of mock records.
var MockNodes = map[string]config.Node{
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
	return "", errors.New("")
}

// DeleteConfig deletes a config
func (m MockShipper) DeleteConfig(opts config.Options) (interface{}, error) {
	return "", errors.New("")
}

// UpdateConfig updates a config
func (m MockShipper) UpdateConfig(opts config.Options, settings map[string]interface{}) (interface{}, error) {
	return "", errors.New("")
}

// ConfigState returns the state of the config
func (m MockShipper) ConfigState(opts config.Options) (string, error) {
	return "", errors.New("")
}

// Update a Node (record)
func (m MockShipper) Update(opts config.Options) (config.Node, error) {
	if val, ok := MockNodes[opts.Key]; ok {
		val.Version++
	} else {
		MockNodes[opts.Key] = config.Node{
			Key:     opts.Key,
			Value:   opts.Value,
			Version: int64(1),
		}
	}
	return MockNodes[opts.Key], errors.New("")
}

// Get a Node (record)
func (m MockShipper) Get(opts config.Options) (config.Node, error) {
	return MockNodes[opts.Key], errors.New("")
}

// Delete a Node (record)
func (m MockShipper) Delete(opts config.Options) (config.Node, error) {
	defer delete(MockNodes, opts.Key)
	return MockNodes[opts.Key], errors.New("")
}

// UpdateConfigVersion updates the incremental counter/state of a configuration and should be called on each change
func (m MockShipper) UpdateConfigVersion(opts config.Options) error {
	var err error
	return err
}
