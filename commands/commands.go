// Command functions
package commands

import (
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/storage"
	"io/ioutil"
	//"log"
)

// Creates a new configuration
func CreateCfg(Config config.Config, args []string) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "create",
	}
	if len(args) > 0 {
		success, _, err := storage.CreateConfig(Config, args[0])
		if err != nil {
			resp.Error = err.Error()
		}
		if success {
			resp.Message = "Successfully created the configuration"
		}
	} else {
		resp.Error = NotEnoughArgsMsg
		// TODO: Error code for this, message may not be necessary - is it worthwhile to try and figure out exactly which arguments were missing?
		// Maybe a future thing to do. I need to git er done right now.
	}
	return resp
}

// Sets a discfg configuration to use for all future commands until unset (it is optional, but conveniently saves a CLI argument - kinda like MongoDB's use)
func Use(Config config.Config, args []string) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "use",
	}
	if len(args) > 0 {
		cc := []byte(args[0])
		err := ioutil.WriteFile(".discfg", cc, 0644)
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Message = "Set current working discfg to " + args[0]
			resp.CurrentDiscfg = args[0]
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}

// Shows which discfg configuration is currently active for use
func Which(Config config.Config, args []string) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "which",
	}
	currentCfg := getDiscfgNameFromFile()
	if currentCfg == "" {
		resp.Error = "No current working configuration has been set at this path."
	} else {
		resp.Message = "Current working configuration: " + currentCfg
		resp.CurrentDiscfg = currentCfg
	}
	return resp
}

// Sets a key value for a given configuration
func SetKey(Config config.Config, args []string) config.ResponseObject {
	resp := config.ResponseObject{
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
			resp.Node.Key = key
			resp.Node.Value = []byte(value)
			resp.Node.Version = 1

			// Only set PrevNode if there was a previous value
			if storageResponse.Value != nil {
				resp.PrevNode = storageResponse
				resp.PrevNode.Key = key
				// Update the current node's value if there was a previous version
				resp.Node.Version = resp.PrevNode.Version + 1
			}
		}
	} else {
		resp.Error = NotEnoughArgsMsg
		if keyErr != nil {
			resp.Error = keyErr.Error()
		}
	}
	return resp
}

func GetKey(Config config.Config, args []string) config.ResponseObject {
	resp := config.ResponseObject{
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
			resp.Error = err.Error()
		}
		if success {
			resp.Node = storageResponse
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}

func DeleteKey(Config config.Config, args []string) config.ResponseObject {
	resp := config.ResponseObject{
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
			resp.Node = storageResponse
			resp.Node.Key = key
			resp.Node.Value = nil
			resp.Node.Version = storageResponse.Version + 1
			resp.PrevNode.Key = key
			resp.PrevNode.Version = storageResponse.Version
			resp.PrevNode.Value = storageResponse.Value
			// log.Println(storageResponse)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}
