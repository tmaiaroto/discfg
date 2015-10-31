package storage

import (
	//"encoding/json"
	"errors"
	"github.com/sdming/gosnow"
	"github.com/tmaiaroto/discfg/config"
	ddb "github.com/tmaiaroto/discfg/storage/dynamodb"
	//"log"
	//"net/http"
)

// A shipper can send information into a database or log etc. While DynamoDB is the planned data store,
// who knows what will happen in the future. A simple interface never hurts.
type Shipper interface {
	CreateConfig(config.Config, string) (bool, interface{}, error)
	Update(config.Config, string, string, string) (bool, interface{}, error)
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

// Generates a new Snowflake
func generateId() (uint64, error) {
	// TODO: return the error so we can do something. Maybe this function isn't even needed...
	v, _ := gosnow.Default()

	// Alternatively you can set the worker id if you are running multiple snowflakes
	// TODO
	// v, err := gosnow.NewSnowFlake(100)

	return v.Next()
}

// Creates a new configuration returning success true/false along with any response and error.
func CreateConfig(cfg config.Config, name string) (bool, interface{}, error) {
	var err error
	if s, ok := shippers[cfg.StorageInterfaceName]; ok {
		return s.CreateConfig(cfg, name)
	} else {
		err = errors.New("Invalid shipper adapter.")
	}
	return false, nil, err
}

// Updates a key value in the configuration
func Update(cfg config.Config, name string, key string, value string) (bool, interface{}, error) {
	var err error
	if s, ok := shippers[cfg.StorageInterfaceName]; ok {
		return s.Update(cfg, name, key, value)
	} else {
		err = errors.New("Invalid shipper adapter.")
	}
	return false, nil, err
}
