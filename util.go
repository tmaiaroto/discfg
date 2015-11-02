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

// {
//     "action": "set",
//     "node": {
//         "createdIndex": 2,
//         "key": "/message",
//         "modifiedIndex": 2,
//         "value": "Hello world"
//     }
// }

// Unlike etcd, we call our "index" an "id" ... But (for now) a history will not be kept.
// These are not atomic counters (can't be, DynamoDB is distributed). They are snowflakes.
// Snowflake generates "roughly" sortable values. Might get interesting. Subject to change.
// I'd love to call it an index, but that would be misleading and technically inaccurate.
// However, don't think the snowflake is useless. It serves another important purpose
// when it comes to DynamoDB. It helps distribute the data.
// See: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GuidelinesForTables.html#GuidelinesForTables.UniformWorkload
// And: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html#WorkingWithTables.primary.key
// TODO: Look into DynamoDB's StreamSpecification ... can help with previous values perhaps saving some work.
type Node struct {
	Id    uint64      `json:"id,omitempty"`
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// Regardless, we will still return what the previous value was just like etcd.
// Even if there's no way to ever return to that value...At least not from discfg (for now).
type PrevNode struct {
	Id    uint64      `json:"id,omitempty"`
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
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
