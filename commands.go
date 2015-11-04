// Command functions
package main

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/tmaiaroto/discfg/storage"
	"io/ioutil"
	"log"
	"strconv"
)

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

	key, keyErr := formatKeyName(key)
	if enoughArgs && keyErr == nil {
		success, storageResponse, err := storage.Update(Config, discfgName, key, value)
		if err != nil {
			resp.Error = "Error updating key value"
			resp.Message = err.Error()
		}
		if success {
			resp.Success = "Successfully updated key value"
			resp.Node.Key = key
			resp.Node.Value = value
			resp.Node.Version = 1

			// Only set PrevNode if there was a previous value
			r := storageResponse.(map[string]string)
			if val, ok := r["value"]; ok {
				resp.PrevNode.Key = key
				resp.PrevNode.Value = val
				prevVersion, _ := strconv.ParseInt(r["version"], 10, 64)
				resp.PrevNode.Version = prevVersion
				// Update the current node's value if there was a previous version
				resp.Node.Version = prevVersion + 1
			}
		}
	} else {
		resp.Error = NotEnoughArgsMsg
		if keyErr != nil {
			resp.Error = keyErr.Error()
		}
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

	key, keyErr := formatKeyName(key)
	if enoughArgs && keyErr == nil {
		success, storageResponse, err := storage.Get(Config, discfgName, key)
		if err != nil {
			resp.Error = "Error getting key value"
			resp.Message = err.Error()
		}
		if success {
			// TODO: refactor. use the types so stroage.Get() returns the type.
			// it would be much nicer.
			r := storageResponse.(map[string]string)
			//parsedId, _ := strconv.ParseUint(r["id"], 10, 64)
			//resp.Node.Id = parsedId

			parsedVersion, _ := strconv.ParseInt(r["version"], 10, 64)
			resp.Node.Version = parsedVersion
			resp.Node.Key = key
			resp.Node.Value = r["value"]
			log.Println(isJSON(resp.Node.Value))
			if isJSON(resp.Node.Value) {
				resp.Node.Raw = json.RawMessage(r["value"])
				//resp.Node.Value = ""
			}

			// log.Println(storageResponse)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	out(resp)
}

func deleteKey(cmd *cobra.Command, args []string) {
	resp := ReponseObject{
		Action: "delete",
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

	key, keyErr := formatKeyName(key)
	if enoughArgs && keyErr == nil {
		success, storageResponse, err := storage.Delete(Config, discfgName, key)
		if err != nil {
			resp.Error = "Error getting key value"
			resp.Message = err.Error()
		}
		if success {
			r := storageResponse.(map[string]string)
			//parsedId, _ := strconv.ParseUint(r["id"], 10, 64)
			//resp.Node.Id = parsedId

			parsedVersion, _ := strconv.ParseInt(r["version"], 10, 64)
			resp.Node.Version = (parsedVersion + 1)
			resp.Node.Key = key
			resp.PrevNode.Version = parsedVersion
			resp.PrevNode.Key = key
			resp.PrevNode.Value = r["value"]
			// log.Println(storageResponse)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	out(resp)
}
