package storage

import (
	//"encoding/json"
	"errors"
	"github.com/tmaiaroto/discfg/config"
	ddb "github.com/tmaiaroto/discfg/storage/dynamodb"
	//"log"
	//"net/http"
)

// A shipper can send information into a database or log etc. While DynamoDB is the planned data store,
// who knows what will happen in the future. A simple interface never hurts.
type Shipper interface {
	CreateConfig(config.Options) (bool, interface{}, error)
	Update(config.Options) (bool, config.Node, error)
	Get(config.Options) (bool, config.Node, error)
	Delete(config.Options) (bool, config.Node, error)
	UpdateConfigVersion(config.Options) bool
}

// Standard shipper result contains errors and other information.
type ShipperResult struct {
	Interface string `json:"interface"`
	Error     error  `json:"error"`
}

// Defines all of the available shipper interfaces available for use.
var shippers = map[string]Shipper{
	"dynamodb": ddb.DynamoDB{},
}

// Creates a new configuration returning success true/false along with any response and error.
func CreateConfig(opts config.Options) (bool, interface{}, error) {
	var err error
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.CreateConfig(opts)
	} else {
		err = errors.New("Invalid shipper adapter.")
	}
	return false, nil, err
}

// Updates a key value in the configuration
func Update(opts config.Options) (bool, config.Node, error) {
	var err error
	var node config.Node
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		if UpdateConfigVersion(opts) {
			return s.Update(opts)
		}
	} else {
		err = errors.New("Invalid shipper adapter.")
	}
	return false, node, err
}

// Gets a key value in the configuration
func Get(opts config.Options) (bool, config.Node, error) {
	var err error
	var node config.Node
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.Get(opts)
	} else {
		err = errors.New("Invalid shipper adapter.")
	}
	return false, node, err
}

// Deletes a key value in the configuration
func Delete(opts config.Options) (bool, config.Node, error) {
	var err error
	var node config.Node
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		if UpdateConfigVersion(opts) {
			return s.Delete(opts)
		}

	} else {
		err = errors.New("Invalid shipper adapter.")
	}
	return false, node, err
}

// Updates the global discfg config version and modified timestamp (on the root key "/")
func UpdateConfigVersion(opts config.Options) bool {
	// Technically, this modified timestamp won't be accurate. The config would have changed already by this point.
	// TODO: Perhaps pass a timestamp to this function to get a little closer
	if s, ok := shippers[opts.StorageInterfaceName]; ok {
		return s.UpdateConfigVersion(opts)
	}
	return false
}
