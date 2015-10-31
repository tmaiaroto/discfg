package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tmaiaroto/discfg/storage"
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
			resp.Success = "Successfully created the configuration."
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
	currentCfg, err := ioutil.ReadFile(".discfg")
	if err != nil {
		resp.Message = "No current working configuration has been set at this path."
	} else {
		resp.Message = "Current working configuration: " + string(currentCfg)
		resp.CurrentDiscfg = string(currentCfg)
	}
	out(resp)
}

// Sets a key value for a given configuration
func setKey(cmd *cobra.Command, args []string) {
	resp := ReponseObject{
		Action: "set",
	}
	if len(args) > 1 {
		success, _, err := storage.Update(Config, args[0], args[1], args[2])
		if err != nil {
			resp.Error = "Error Creating Configuration"
			resp.Message = err.Error()
		}
		if success {
			resp.Success = "Successfully created the configuration."
			// TODO: a verbose, vv, or debug mode which would include the response from AWS
			// So if verbose, then Message would take on this response...Or perhaps another field.
			//log.Println(response)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
}
