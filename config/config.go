// Defines various structures including the configuration.
package config

import (
	"encoding/json"
)

// All of the configuration options needed by various functions (set by CLI commands, config files, other code, etc.)
type Options struct {
	CfgName              string
	ConditionalValue     string
	Recursive            bool
	Key                  string
	Value                string
	TTL                  int64
	StorageInterfaceName string
	Storage              struct {
		DynamoDB struct {
			Region             string
			AccessKeyId        string
			SecretAccessKey    string
			CredProfile        string
			WriteCapacityUnits int64
			ReadCapacityUnits  int64
		}
	}
	Version      string
	OutputFormat string
}

// The response for output.
type ResponseObject struct {
	Action string `json:"action"`
	//Node          Node     `json:"node,omitempty"`
	Node          Node   `json:"node,omitempty"`
	PrevNode      Node   `json:"prevNode,omitempty"`
	ErrorCode     int    `json:"errorCode,omitempty"`
	CurrentDiscfg string `json:"currentDiscfg,omitempty"`
	Error         string `json:"error,omitempty"`
	Message       string `json:"message,omitempty"`
	// Information about the config
	CfgVersion int64 `json:"cfgVersion,omitempty"`
	// In seconds since that's probably more common for people
	CfgModified int64 `json:"cfgModified,omitempty"`
	// In nanoseconds for the gophers like me who are snobby about time =)
	CfgModifiedNanoseconds int64 `json:"cfgModifiedNanoseconds,omitempty"`
	// A parsed date for humans to read
	CfgModifiedParsed string `json:"cfgModifiedParsed,omitempty"`
}

// NOTES ON NODES:
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
// So with that in mind, a simple version is found on each node. While a bit naive, it's effective for many situations.
// Not seen on this struct (for now), but stored in DynamoDB is also a list of the parent nodes (full paths).
// This is for traversing needs.
//
// Another great piece of DynamoDB documentation with regard to counters and conditional writes can be found here:
// http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.AtomicCounters
//
// Again, this highlights discfg's source of inspriation (etcd) and difference from it.
//
// TODO: Content-Type?
//
type Node struct {
	Version                int64           `json:"version,omitempty"`
	Key                    string          `json:"key,omitempty"`
	Value                  []byte          `json:"-"` //`json:"value,omitempty"`
	OutputValue            json.RawMessage `json:"value,omitempty"`
	Nodes                  []Node          `json:"nodes,omitepty"`
	CfgVersion             int64           `json:"-"`
	CfgModifiedNanoseconds int64           `json:"-"`
}
