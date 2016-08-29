// Package config defines various structures including the configuration.
package config

import (
	//"encoding/json"
	"time"
)

// Options needed by various functions (set by CLI commands, config files, other code, etc.)
type Options struct {
	CfgName              string
	ConditionalValue     string
	Recursive            bool
	Key                  string
	Value                []byte
	TTL                  int64
	StorageInterfaceName string
	// Storage options, for now AWS is the only supported storage
	Storage struct {
		AWS
	}
	Version      string
	OutputFormat string
}

// AWS credentials and options
type AWS struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	CredProfile     string
}

// ResponseObject for output
type ResponseObject struct {
	Action        string `json:"action"`
	Item          Item   `json:"item,omitempty"`
	PrevItem      Item   `json:"prevItem,omitempty"`
	ErrorCode     int    `json:"errorCode,omitempty"`
	CurrentDiscfg string `json:"currentDiscfg,omitempty"`
	// Error message
	Error string `json:"error,omitempty"`
	// Message returned to the CLI
	Message string `json:"message,omitempty"`
	// Add this? Might be useful for troubleshooting, but users shouldn't really need to worry about it.
	// On the other hand, for things like DynamoDB, it's handy to know where the config stands in terms of scalability/capacity.
	// For things like S3, there's no real settings to care about (for the most part at least).
	// StorageResponse interface{} `json:"storageResponse,omitempty"`
	// Information about the config
	CfgVersion int64 `json:"cfgVersion,omitempty"`
	// In seconds since that's probably more common for people
	CfgModified int64 `json:"cfgModified,omitempty"`
	// In nanoseconds for the gophers like me who are snobby about time =)
	CfgModifiedNanoseconds int64 `json:"cfgModifiedNanoseconds,omitempty"`
	// A parsed date for humans to read
	CfgModifiedParsed string `json:"cfgModifiedParsed,omitempty"`
	// Configuration state (some storage engines, such as DynamoDB, have "active" and "updating" states)
	CfgState string `json:"cfgState,omitempty"`
	// Information about the configuration storage
	CfgStorage StorageInfo `json:"cfgStorage,omitempty"`
}

// StorageInfo holds information about the storage engine used for the configuration
type StorageInfo struct {
	Name          string                 `json:"name"`
	InterfaceName string                 `json:"interfaceName"`
	Options       map[string]interface{} `json:"options"`
}

// NOTES ON ITEMS (somewhat similar to etcd's nodes):
// Unlike etcd, there is no "index" key because discfg doesn't try to be a state machine like etcd.
// The index there refers to some internal state of the entire system and certain actions advance that state.
// discfg does not have a distributed lock system nor this sense of global state.
//
// However, it is useful for applications (and humans) to get a sense of change. So two thoughts:
//   1. An "id" value using snowflake (so it's roughly sortable - the thought being sequential enough for discfg's needs)
//   2. A "version" value that simply increments on each update
//
// If using snowflake ids, it would make sense to add those values as part of the index (RANGE). It would certainly
// help DynaoDB distribute the data...
// See: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GuidelinesForTables.html#GuidelinesForTables.UniformWorkloa
// And: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html#WorkingWithTables.primary.key
//
// The challenge then here is there wouldn't be conditional updates by DynamoDB design. Those would need to be added
// and it would require more queries. The database would be append only (which has its own benefits). Then there would
// eventually need to be some sort of expiration on old items. Since one of the gaols of discfg is cost efficiency,
// it doesn't make sense to keep old items around. Plus, going backwards in time is not a typical need for a configuration
// service. The great thing about etcd's state here is the ability to watch for changes and should that HTTP connection
// be interrupted, it could be resumed from a specific point. This is just one reason for that state index.
//
// discfg does not have this feature. There is no way to watch for a key update because discfg is not meant to run in
// persistence. The data is of course, but the service is not. It's designed to run on demand CLI or AWS Lambda.
// It's simply a different design decision in order to hit a goal. discfg's answer for this need would be to reach for
// other AWS services to push notifications out (SNS), add to a message queue (SQS), etc.
//
// So with that in mind, a simple version is found on each item. While a bit naive, it's effective for many situations.
// Not seen on this struct (for now), but stored in DynamoDB is also a list of the parent items (full paths).
// This is for traversing needs.
//
// Another great piece of DynamoDB documentation with regard to counters and conditional writes can be found here:
// http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.AtomicCounters
//
// Again, this highlights discfg's source of inspriation (etcd) and difference from it.
//
// TODO: Content-Type?
// Value will be an interface{}, but stored as []byte in DynamoDB.
// Other storage engines may convert to something else.
// For now, all data is coming in as string. Either from the terminal or a RESTful API.

// Item defines the data structure around a key and its state
type Item struct {
	Version int64  `json:"version,omitempty"`
	Key     string `json:"key,omitempty"`
	//Value   []byte `json:"value,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	OutputValue map[string]interface{} `json:"ovalue,omitempty"`

	// perfect for json, not good if some other value was stored
	//OutputValue            map[string]interface{} `json:"value,omitempty"`
	// We really need interface{} for any type of data. []byte above is for DynamoDB specifically.
	// It could be ... yea. an interface{} too. converted to []byte for storing in dynamodb.
	//OutputValue            interface{} `json:"value,omitempty"`
	TTL              int64     `json:"ttl,omitempty"`
	Expiration       time.Time `json:"-"`
	OutputExpiration string    `json:"expiration,omitempty"`
	// For now, skip this. The original thinking was to have a tree like directory structure like etcd.
	// Though discfg has now deviated away from that to a flat key/value structure.
	// Items                  []Item    `json:"items,omitepty"`
	CfgVersion             int64 `json:"-"`
	CfgModifiedNanoseconds int64 `json:"-"`
}
