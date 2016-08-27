// Package commands contains functions to work with a discfg configuration from a high level.
package commands

import (
	"bytes"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/storage"
	"io/ioutil"
	"strconv"
	"time"
)

// CreateCfg creates a new configuration
func CreateCfg(opts config.Options, settings map[string]interface{}) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "create cfg",
	}
	if len(opts.CfgName) > 0 {
		_, err := storage.CreateConfig(opts, settings)
		if err != nil {
			resp.Error = err.Error()
			resp.Message = "Error creating the configuration"
		} else {
			resp.Message = "Successfully created the configuration"
		}
	} else {
		resp.Error = NotEnoughArgsMsg
		// TODO: Error code for this, message may not be necessary - is it worthwhile to try and figure out exactly which arguments were missing?
		// Maybe a future thing to do. I need to git er done right now.
	}
	return resp
}

// DeleteCfg deletes a configuration
func DeleteCfg(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "delete cfg",
	}
	if len(opts.CfgName) > 0 {
		_, err := storage.DeleteConfig(opts)
		if err != nil {
			resp.Error = err.Error()
			resp.Message = "Error deleting the configuration"
		} else {
			resp.Message = "Successfully deleted the configuration"
		}
	} else {
		resp.Error = NotEnoughArgsMsg
		// TODO: Error code for this, message may not be necessary - is it worthwhile to try and figure out exactly which arguments were missing?
		// Maybe a future thing to do. I need to git er done right now.
	}
	return resp
}

// UpdateCfg updates a configuration's options/settings (if applicable, depends on the interface)
func UpdateCfg(opts config.Options, settings map[string]interface{}) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "update cfg",
	}

	// Note: For some storage engines, such as DynamoDB, it could take a while for changes to be reflected.
	if len(settings) > 0 {
		_, updateErr := storage.UpdateConfig(opts, settings)
		if updateErr != nil {
			resp.Error = updateErr.Error()
			resp.Message = "Error updating the configuration"
		} else {
			resp.Message = "Successfully updated the configuration"
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}

	return resp
}

// Use sets a discfg configuration to use for all future commands until unset (it is optional, but conveniently saves a CLI argument - kinda like MongoDB's use)
func Use(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "use",
	}
	if len(opts.CfgName) > 0 {
		cc := []byte(opts.CfgName)
		err := ioutil.WriteFile(".discfg", cc, 0644)
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Message = "Set current working discfg to " + opts.CfgName
			resp.CurrentDiscfg = opts.CfgName
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}

// Which shows which discfg configuration is currently active for use
func Which(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "which",
	}
	currentCfg := GetDiscfgNameFromFile()
	if currentCfg == "" {
		resp.Error = NoCurrentWorkingCfgMsg
	} else {
		resp.Message = "Current working configuration: " + currentCfg
		resp.CurrentDiscfg = currentCfg
	}
	return resp
}

// SetKey sets a key value for a given configuration
func SetKey(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "set",
	}
	// Do not allow empty values to be set
	if opts.Value == nil {
		resp.Error = ValueRequiredMsg
		return resp
	}

	if opts.CfgName == "" {
		resp.Error = MissingCfgNameMsg
		return resp
	}

	key, keyErr := formatKeyName(opts.Key)
	if keyErr == nil {
		opts.Key = key
		storageResponse, err := storage.Update(opts)
		if err != nil {
			resp.Error = err.Error()
			resp.Message = "Error updating key value"
		} else {
			resp.Item.Key = key
			resp.Item.Value = opts.Value
			resp.Item.Version = 1

			// Only set PrevItem if there was a previous value
			if storageResponse.Value != nil {
				resp.PrevItem = storageResponse
				resp.PrevItem.Key = key
				// Update the current item's value if there was a previous version
				resp.Item.Version = resp.PrevItem.Version + 1
			}
		}
	} else {
		resp.Error = keyErr.Error()
	}
	return resp
}

// GetKey gets a key from a configuration
func GetKey(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "get",
	}
	key, keyErr := formatKeyName(opts.Key)
	if keyErr == nil {
		opts.Key = key
		storageResponse, err := storage.Get(opts)
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Item = storageResponse
		}
	} else {
		resp.Error = keyErr.Error()
	}
	return resp
}

// DeleteKey deletes a key from a configuration
func DeleteKey(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "delete",
	}
	key, keyErr := formatKeyName(opts.Key)
	if keyErr == nil {
		opts.Key = key
		storageResponse, err := storage.Delete(opts)
		if err != nil {
			resp.Error = err.Error()
			resp.Message = "Error getting key value"
		} else {
			resp.Item = storageResponse
			resp.Item.Key = opts.Key
			resp.Item.Value = nil
			resp.Item.Version = storageResponse.Version + 1
			resp.PrevItem.Key = opts.Key
			resp.PrevItem.Version = storageResponse.Version
			resp.PrevItem.Value = storageResponse.Value
			// log.Println(storageResponse)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}

// Info about the configuration including global version/state and modified time
func Info(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "info",
	}

	if opts.CfgName != "" {
		// Just get the root key
		opts.Key = "/"

		storageResponse, err := storage.Get(opts)
		if err != nil {
			resp.Error = err.Error()
		} else {
			// Debating putting the item value on here... (allowing users to store values on the config or "root")
			// resp.Item = storageResponse
			// Set the configuration version and modified time on the response
			// Item.CfgVersion and Item.CfgModifiedNanoseconds are not included in the JSON output
			resp.CfgVersion = storageResponse.CfgVersion
			resp.CfgModified = 0
			resp.CfgModifiedNanoseconds = storageResponse.CfgModifiedNanoseconds
			// Modified in seconds
			resp.CfgModified = storageResponse.CfgModifiedNanoseconds / int64(time.Second)
			// Modified parsed
			modified := time.Unix(0, storageResponse.CfgModifiedNanoseconds)
			resp.CfgModifiedParsed = modified.Format(time.RFC3339)

			// Get the status (only applicable for some storage interfaces, such as DynamoDB)
			resp.CfgState, err = storage.ConfigState(opts)
			if err != nil {
				resp.Error = err.Error()
			} else {
				var buffer bytes.Buffer
				buffer.WriteString(opts.CfgName)
				if resp.CfgState != "" {
					buffer.WriteString(" (")
					buffer.WriteString(resp.CfgState)
					buffer.WriteString(")")
				}
				buffer.WriteString(" version ")
				buffer.WriteString(strconv.FormatInt(resp.CfgVersion, 10))
				buffer.WriteString(" last modified ")
				buffer.WriteString(modified.Format(time.RFC1123))
				resp.Message = buffer.String()
				buffer.Reset()
			}
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}

// Export a discfg to file in JSON format
func Export(opts config.Options, args []string) {
	// TODO
}
