package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tmaiaroto/discfg/storage"
	"io/ioutil"
	// "log"
	"strconv"
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

// Creates a new configuration
func createCfg(cmd *cobra.Command, args []string) {
	resp := ReponseObject{
		Action: "create",
	}
	if len(args) > 0 {
		success, _, err := storage.CreateConfig(Config, args[0])
		if err != nil {
			resp.Error = "Error Creating Configuration"
			resp.Message = err.Error()
		}
		if success {
			resp.Success = "Successfully created the configuration"
			// TODO: a verbose, vv, or debug mode which would include the response from AWS
			// So if verbose, then Message would take on this response...Or perhaps another field.
			//log.Println(response)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
		// TODO: Error code for this, message may not be necessary - is it worthwhile to try and figure out exactly which arguments were missing?
		// Maybe a future thing to do. I need to git er done right now.
	}
	out(resp)
}

// Sets a discfg configuration to use for all future commands until unset (it is optional, but conveniently saves a CLI argument - kinda like MongoDB's use)
func use(cmd *cobra.Command, args []string) {
	resp := ReponseObject{
		Action: "use",
	}
	if len(args) > 0 {
		cc := []byte(args[0])
		err := ioutil.WriteFile(".discfg", cc, 0644)
		if err != nil {
			resp.Error = "There was a problem setting the discfg to use"
			resp.Message = err.Error()
		} else {
			resp.Success = "Set current working discfg to " + args[0]
			resp.CurrentDiscfg = args[0]
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	out(resp)
}

// Shows which discfg configuration is currently active for use
func which(cmd *cobra.Command, args []string) {
	resp := ReponseObject{
		Action: "which",
	}
	currentCfg := getDiscfgNameFromFile()
	if currentCfg != "" {
		resp.Message = "No current working configuration has been set at this path."
	} else {
		resp.Message = "Current working configuration: " + currentCfg
		resp.CurrentDiscfg = currentCfg
	}
	out(resp)
}

// Just returns the name of the set discfg name (TODO: will need to change as .discfg gets more complex)
func getDiscfgNameFromFile() string {
	currentCfg, err := ioutil.ReadFile(DiscfgFileName)
	if err == nil {
		return string(currentCfg)
	}
	return ""
}

// Sets a key value for a given configuration
func setKey(cmd *cobra.Command, args []string) {
	resp := ReponseObject{
		Action: "set",
	}
	// TODO: refactor
	var discfgName string
	var key string
	var value string
	enoughArgs := false
	if len(args) > 1 {
		currentName := getDiscfgNameFromFile()
		if len(args) == 2 && currentName != "" {
			discfgName = currentName
			key = args[0]
			value = args[1]
			enoughArgs = true
		} else {
			if len(args) == 3 {
				discfgName = args[0]
				key = args[1]
				value = args[2]
				enoughArgs = true
			}
		}
	}

	if enoughArgs {
		success, storageResponse, err := storage.Update(Config, discfgName, key, value)
		if err != nil {
			resp.Error = "Error updating key value"
			resp.Message = err.Error()
		}
		if success {
			resp.Success = "Successfully updated key value"
			resp.Node.Key = key
			resp.Node.Value = value

			resp.Message = storageResponse.(string)

			// TODO: a verbose, vv, or debug mode which would include the response from AWS
			// So if verbose, then Message would take on this response...Or perhaps another field.
			//log.Println(response)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	out(resp)
}

func getKey(cmd *cobra.Command, args []string) {
	resp := ReponseObject{
		Action: "get",
	}
	// TODO: refactor
	var discfgName string
	var key string
	enoughArgs := false
	if len(args) > 0 {
		currentName := getDiscfgNameFromFile()
		if len(args) == 1 && currentName != "" {
			discfgName = currentName
			key = args[0]
			enoughArgs = true
		} else {
			if len(args) == 2 {
				discfgName = args[0]
				key = args[1]
				enoughArgs = true
			}
		}
	}

	if enoughArgs {
		success, storageResponse, err := storage.Get(Config, discfgName, key)
		if err != nil {
			resp.Error = "Error getting key value"
			resp.Message = err.Error()
		}
		if success {
			// TODO: refactor. use the types so stroage.Get() returns the type.
			// it would be much nicer.
			r := storageResponse.(map[string]string)
			parsedId, _ := strconv.ParseUint(r["id"], 10, 64)
			resp.Node.Id = parsedId
			resp.Node.Key = key
			resp.Node.Value = r["value"]
			// log.Println(storageResponse)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	out(resp)
}
