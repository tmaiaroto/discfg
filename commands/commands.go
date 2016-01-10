// Command functions
package commands

import (
	"bytes"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/storage"
	"io/ioutil"
	"strconv"
	"time"
	// "log"
)

// Creates a new configuration
func CreateCfg(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "create",
	}
	if len(opts.CfgName) > 0 {
		success, _, err := storage.CreateConfig(opts)
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

// Shows which discfg configuration is currently active for use
func Which(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "which",
	}
	currentCfg := GetDiscfgNameFromFile()
	if currentCfg == "" {
		resp.Error = "No current working configuration has been set at this path."
	} else {
		resp.Message = "Current working configuration: " + currentCfg
		resp.CurrentDiscfg = currentCfg
	}
	return resp
}

// Sets a key value for a given configuration
func SetKey(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "set",
	}

	key, keyErr := formatKeyName(opts.Key)
	if keyErr == nil {
		opts.Key = key
		success, storageResponse, err := storage.Update(opts)
		if err != nil {
			resp.Error = "Error updating key value"
			resp.Message = err.Error()
		}
		if success {
			resp.Node.Key = key
			resp.Node.Value = opts.Value //[]byte(opts.Value)
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

// Gets a key
func GetKey(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "get",
	}
	key, keyErr := formatKeyName(opts.Key)
	if keyErr == nil {
		opts.Key = key
		success, storageResponse, err := storage.Get(opts)
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

// Deletes a key
func DeleteKey(opts config.Options) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "delete",
	}
	key, keyErr := formatKeyName(opts.Key)
	if keyErr == nil {
		opts.Key = key
		success, storageResponse, err := storage.Delete(opts)
		if err != nil {
			resp.Error = "Error getting key value"
			resp.Message = err.Error()
		}
		if success {
			resp.Node = storageResponse
			resp.Node.Key = opts.Key
			resp.Node.Value = nil
			resp.Node.Version = storageResponse.Version + 1
			resp.PrevNode.Key = opts.Key
			resp.PrevNode.Version = storageResponse.Version
			resp.PrevNode.Value = storageResponse.Value
			// log.Println(storageResponse)
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}

// Information about the configuration including global version and modified time
func Info(opts config.Options, args []string) config.ResponseObject {
	resp := config.ResponseObject{
		Action: "info",
	}

	if opts.CfgName != "" {
		// Just get the root key
		opts.Key = "/"
		success, storageResponse, err := storage.Get(opts)
		if err != nil {
			resp.Error = err.Error()
		}
		if success {
			resp.Node = storageResponse
			// Set the configuration version and modified time on the response
			// Node.CfgVersion and Node.CfgModifiedNanoseconds are not included in the JSON output
			resp.CfgVersion = resp.Node.CfgVersion
			resp.CfgModified = 0
			resp.CfgModifiedNanoseconds = resp.Node.CfgModifiedNanoseconds
			if resp.Node.CfgModifiedNanoseconds > 0 {
				resp.CfgModified = resp.Node.CfgModifiedNanoseconds / 1000000000
			}
			modified := time.Unix(resp.CfgModified, 0)
			resp.CfgModifiedParsed = modified.Format(time.RFC3339)

			var buffer bytes.Buffer
			buffer.WriteString(opts.CfgName)
			buffer.WriteString(" version ")
			buffer.WriteString(strconv.FormatInt(resp.CfgVersion, 10))
			buffer.WriteString(" last modified ")
			buffer.WriteString(modified.Format(time.RFC1123))
			resp.Message = buffer.String()
			buffer.Reset()
		}
	} else {
		resp.Error = NotEnoughArgsMsg
	}
	return resp
}

// Exports a discfg to file in JSON format
func Export(opts config.Options, args []string) {
	// TODO
}
