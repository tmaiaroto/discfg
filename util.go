// Various utilities used by commands are found in this file as well as response structs, constants, etc.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	//"github.com/asaskevich/govalidator"
	"io/ioutil"
)

// Just like etcd, the response object is built in a very similar manner.
type ReponseObject struct {
	Action        string   `json:"action"`
	Node          Node     `json:"node,omitempty"`
	PrevNode      PrevNode `json:"prevNode,omitempty"`
	ErrorCode     int      `json:"errorCode,omitempty"`
	Message       string   `json:"message,omitempty"`
	CurrentDiscfg string   `json:"currentDiscfg,omitempty"`
	// Error and Success are human readable short messages meant for CLI, not for JSON response.
	Error   string `json:"-"`
	Success string `json:"-"`
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
type Node struct {
	//Id    uint64      `json:"id,omitempty"`
	Version int64       `json:"version,omitempty"`
	Key     string      `json:"key,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

// On an update, the previous node will also be returned.
type PrevNode struct {
	//Id    uint64      `json:"id,omitempty"`
	Version int64       `json:"version,omitempty"`
	Key     string      `json:"key,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

const NotEnoughArgsMsg = "Not enough arguments passed. Run 'discfg help' for usage."
const DiscfgFileName = ".discfg"

// Output
func out(resp ReponseObject) {
	switch Config.OutputFormat {
	case "json":
		o, _ := json.Marshal(resp)
		// TODO: verbose mode here too? Shouldn't be in a situation where it can't be marshaled but who knows.
		// Always best to handle errors.
		// if(oErr) {
		// 	errorLabel("Error")
		// 	fmt.Print(oErr)
		// }
		fmt.Print(string(o))
	case "human":
		// Only gets messages... It doesn't get the details of which key was updated, etc. Just that one was or wasn't updated.
		if resp.Error != "" {
			errorLabel(resp.Error)
		}
		if resp.Success != "" {
			successLabel(resp.Success)
		}
		if resp.Message != "" {
			fmt.Print(resp.Message + "\n")
		}
	}
}

// Just returns the name of the set discfg name (TODO: will need to change as .discfg gets more complex).
func getDiscfgNameFromFile() string {
	currentCfg, err := ioutil.ReadFile(DiscfgFileName)
	if err == nil {
		return string(currentCfg)
	}
	return ""
}

// Simple substring function
func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// Checks and formats the key name. Ensures it beings with a "/" and that it's valid, if not it an empty string is returned (which won't be used).
func formatKeyName(key string) (string, error) {
	var err error
	k := ""
	if len(key) > 0 {
		if substr(key, 0, 1) != "/" {
			var buffer bytes.Buffer
			buffer.WriteString("/")
			buffer.WriteString(key)
			k = buffer.String()
			buffer.Reset()
		}
	} else {
		return "", errors.New("Missing key name")
	}

	// Ensure valid characters
	r, _ := regexp.Compile(`[\w\/\-]+$`)
	if !r.MatchString(k) {
		return "", errors.New("Invalid key name")
	}

	// Remove any trailing slashes (unless there's only one, the root)
	if len(k) > 1 {
		for k[len(k)-1:] == "/" {
			k = k[:len(k)-1]
		}
	}

	return k, err
}
