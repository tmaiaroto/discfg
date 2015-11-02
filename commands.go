// Command functions
package main

import (
	"github.com/spf13/cobra"
	"github.com/tmaiaroto/discfg/storage"
	"io/ioutil"
	// "log"
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

			resp.Message = storageResponse.(string)

			// TODO: a verbose, vv, or debug mode which would include the response from AWS
			// So if verbose, then Message would take on this response...Or perhaps another field.
			//log.Println(response)
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
