// Package storage contains the very important Shipper interface which is responsible for working with storage engines.
package storage

import (
	//"encoding/json"
	"errors"
	"github.com/tmaiaroto/discfg/config"
	ddb "github.com/tmaiaroto/discfg/storage/dynamodb"
	//"log"
	//"net/http"
)

// Shipper can send information into a database or log etc. While DynamoDB is the planned data store,
// who knows what will happen in the future. A simple interface never hurts.
type Shipper interface {
	CreateConfig(config.Options, map[string]interface{}) (bool, interface{}, error)
	DeleteConfig(config.Options) (bool, interface{}, error)
	UpdateConfig(config.Options, map[string]interface{}) (bool, interface{}, error)
	ConfigState(config.Options) string
	Update(config.Options) (bool, config.Node, error)
	Get(config.Options) (bool, config.Node, error)
	Delete(config.Options) (bool, config.Node, error)
	UpdateConfigVersion(config.Options) bool
}

// ShipperResult contains errors and other information.
type ShipperResult struct {
	Interface string `json:"interface"`
	Error     error  `json:"error"`
}

// Defines all of the available shipper interfaces available for use.
var shippers = map[string]Shipper{
	"dynamodb": ddb.DynamoDB{},
}

// CreateConfig creates a new configuration returning success true/false along with any response and error.
func CreateConfig(opts config.Options, settings map[string]interface{}) (bool, interface{}, error) {
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.CreateConfig(opts, settings)
	}
	return false, nil, errors.New("Invalid shipper interface.")
}

// DeleteConfig deletes an existing configuration
func DeleteConfig(opts config.Options) (bool, interface{}, error) {
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.DeleteConfig(opts)
	}
	return false, nil, errors.New("Invalid shipper interface.")
}

// UpdateConfig updates the options/settings for a configuration (may not be implementd by each interface)
func UpdateConfig(opts config.Options, settings map[string]interface{}) (bool, interface{}, error) {
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.UpdateConfig(opts, settings)
	}
	return false, nil, errors.New("Invalid shipper interface.")
}

// ConfigState returns the config state (just a simple string message, could be "ACTIVE" for example)
func ConfigState(opts config.Options) string {
	// TODO: May get more elaborate and have codes for this too, but will probably always have a string message
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.ConfigState(opts)
	}
	return ""
}

// Update a key value in the configuration
func Update(opts config.Options) (bool, config.Node, error) {
	var node config.Node
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		if UpdateConfigVersion(opts) {
			return s.Update(opts)
		}
	}
	return false, node, errors.New("Invalid shipper interface.")
}

// Get a key value in the configuration
func Get(opts config.Options) (bool, config.Node, error) {
	var node config.Node
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.Get(opts)
	}
	return false, node, errors.New("Invalid shipper interface.")
}

// Delete a key value in the configuration
func Delete(opts config.Options) (bool, config.Node, error) {
	var node config.Node
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		if UpdateConfigVersion(opts) {
			return s.Delete(opts)
		}
	}
	return false, node, errors.New("Invalid shipper interface.")
}

// UpdateConfigVersion updates the global discfg config version and modified timestamp (on the root key "/")
func UpdateConfigVersion(opts config.Options) bool {
	// Technically, this modified timestamp won't be accurate. The config would have changed already by this point.
	// TODO: Perhaps pass a timestamp to this function to get a little closer
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.UpdateConfigVersion(opts)
	}
	return false
}
