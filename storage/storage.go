// Package storage contains the very important Shipper interface which is responsible for working with storage engines.
package storage

import (
	"errors"
	"github.com/tmaiaroto/discfg/config"
	ddb "github.com/tmaiaroto/discfg/storage/dynamodb"
)

// Shipper can send information into a database or log etc. While DynamoDB is the planned data store,
// who knows what will happen in the future. A simple interface never hurts.
type Shipper interface {
	CreateConfig(config.Options, map[string]interface{}) (interface{}, error)
	DeleteConfig(config.Options) (interface{}, error)
	UpdateConfig(config.Options, map[string]interface{}) (interface{}, error)
	ConfigState(config.Options) (string, error)
	Update(config.Options) (config.Item, error)
	Get(config.Options) (config.Item, error)
	Delete(config.Options) (config.Item, error)
	UpdateConfigVersion(config.Options) error
}

// ShipperResult contains errors and other information.
type ShipperResult struct {
	Interface string `json:"interface"`
	Error     error  `json:"error"`
}

// Error message constants, reduce repetition.
const (
	errMsgInvalidShipper = "Invalid shipper interface."
)

// A map of all Shipper interfaces available for use (with some defaults).
var shippers = map[string]Shipper{
	"dynamodb": ddb.DynamoDB{},
}

// RegisterShipper allows anyone importing discfg into their own project to register new shippers or overwrite the defaults.
func RegisterShipper(name string, shipper Shipper) {
	shippers[name] = shipper
}

// ListShippers returns the list of available shippers.
func ListShippers() map[string]Shipper {
	return shippers
}

// CreateConfig creates a new configuration returning success true/false along with any response and error.
func CreateConfig(opts config.Options, settings map[string]interface{}) (interface{}, error) {
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.CreateConfig(opts, settings)
	}
	return nil, errors.New(errMsgInvalidShipper)
}

// DeleteConfig deletes an existing configuration
func DeleteConfig(opts config.Options) (interface{}, error) {
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.DeleteConfig(opts)
	}
	return nil, errors.New(errMsgInvalidShipper)
}

// UpdateConfig updates the options/settings for a configuration (may not be implementd by each interface)
func UpdateConfig(opts config.Options, settings map[string]interface{}) (interface{}, error) {
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.UpdateConfig(opts, settings)
	}
	return nil, errors.New(errMsgInvalidShipper)
}

// ConfigState returns the config state (just a simple string message, could be "ACTIVE" for example)
func ConfigState(opts config.Options) (string, error) {
	// TODO: May get more elaborate and have codes for this too, but will probably always have a string message
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.ConfigState(opts)
	}
	return "", errors.New(errMsgInvalidShipper)
}

// Update a key value in the configuration
func Update(opts config.Options) (config.Item, error) {
	var item config.Item
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		err := UpdateConfigVersion(opts)
		if err != nil {
			return item, err
		}
		return s.Update(opts)
	}
	return item, errors.New(errMsgInvalidShipper)
}

// Get a key value in the configuration
func Get(opts config.Options) (config.Item, error) {
	var item config.Item
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.Get(opts)
	}
	return item, errors.New(errMsgInvalidShipper)
}

// Delete a key value in the configuration
func Delete(opts config.Options) (config.Item, error) {
	var item config.Item
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		err := UpdateConfigVersion(opts)
		if err != nil {
			return item, err
		}
		return s.Delete(opts)
	}
	return item, errors.New(errMsgInvalidShipper)
}

// UpdateConfigVersion updates the global discfg config version and modified timestamp (on the root key "/")
func UpdateConfigVersion(opts config.Options) error {
	// Technically, this modified timestamp won't be accurate. The config would have changed already by this point.
	// TODO: Perhaps pass a timestamp to this function to get a little closer
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.UpdateConfigVersion(opts)
	}
	return errors.New(errMsgInvalidShipper)
}
