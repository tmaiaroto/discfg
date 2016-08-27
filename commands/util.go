// Package commands utilities and response structs, constants, etc.
package commands

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	//"strconv"
	//"github.com/pquerna/ffjson/ffjson"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/tmaiaroto/discfg/config"
	"io/ioutil"
	"time"
)

// TODO: Refactor: Change the following ...Msg constants and use config/status.go instead to centralize the error codes and messages.

// NotEnoughArgsMsg defines a message for input validation
const NotEnoughArgsMsg = "Not enough arguments passed. Run 'discfg help' for usage."

// ValueRequired defines a message for input validation
const ValueRequiredMsg = "A value is required. Run 'discfg help' for usage."

// DiscfgFileName defines the temporary filename used to hold the current working config name
const DiscfgFileName = ".discfg"

// NoCurrentWorkingCfgMsg defines a message for an error when a config name can not be found in a .discfg file
const NoCurrentWorkingCfgMsg = "No current working configuration has been set at this path."

// MissingKeyNameMsg defines a message for input validation when a key name was not passed
const MissingKeyNameMsg = "Missing key name"

// InvalidKeyNameMsg defines a message for input validation
const InvalidKeyNameMsg = "Invalid key name"

// MissingCfgNameMsg defines a message for input validation
const MissingCfgNameMsg = "Missing configuration name"

// Out formats a config.ResponseObject for suitable output
func Out(opts config.Options, resp config.ResponseObject) config.ResponseObject {
	// We've stored everything as binary data. But that can be many things.
	// A string, a number, or even JSON. We can check to see if it's something we can marshal to JSON.
	// If that fails, then we'll just return it as a string in the JSON response under the "value" key.
	//
	// If it isn't JSON, then return a base64 string.
	// TODO: Add Content-Type field of some sort so there's some context?
	//
	// TODO: Stuff like this will now be handled by an output interface.
	// ...and will also handle the content-type situation.
	// Output JSON, output Msgpack, output Protobuf? output whatever Content-Type.
	//
	// if resp.Item.Value != nil {
	// 	if !isJSON(string(resp.Item.Value)) {
	// 		// Return base64 when not JSON?
	// 		// b64Str := base64.StdEncoding.EncodeToString(resp.Item.Value)
	// 		//resp.Item.Value = []byte(strconv.Quote(b64Str))
	// 		resp.Item.Value = []byte(strconv.Quote(string(resp.Item.Value)))
	// 	}
	// 	// The output value is always raw JSON. It is not stored in the data store.
	// 	// It's simply for display.
	// 	resp.Item.OutputValue = json.RawMessage(resp.Item.Value)
	// }

	// // Same for the PrevItem if set
	// if resp.PrevItem.Value != nil {
	// 	if !isJSON(string(resp.PrevItem.Value)) {
	// 		resp.PrevItem.Value = []byte(strconv.Quote(string(resp.PrevItem.Value)))
	// 	}
	// 	resp.PrevItem.OutputValue = json.RawMessage(resp.PrevItem.Value)
	// }

	// Format the expiration time (if applicable). This prevents output like "0001-01-01T00:00:00Z" when empty
	// and allows for the time.RFC3339Nano format to be used whereas time.Time normally marshals to a different format.
	if resp.Item.TTL > 0 {
		resp.Item.OutputExpiration = resp.Item.Expiration.Format(time.RFC3339Nano)
	}

	switch opts.OutputFormat {
	case "json":
		o, _ := json.Marshal(&resp)
		// TODO: Benchmark this - is it faster?
		// o, _ := ffjson.Marshal(&resp)
		//
		// TODO: verbose mode here too? Shouldn't be in a situation where it can't be marshaled but who knows.
		// Always best to handle errors.
		// if(oErr) {
		// 	errorLabel("Error")
		// 	fmt.Print(oErr)
		// }
		fmt.Print(string(o))
	case "human":
		if resp.Error != "" {
			errorLabel(resp.Error)
		}
		if resp.Item.Value != nil {
			// The value should be a byte array, for th CLI we want a string.
			fmt.Println(string(resp.Item.Value.([]byte)))
		} else {
			if resp.Message != "" {
				fmt.Println(resp.Message)
			}
		}
	}
	return resp
}

// Changes the color for error messages. Good for one line heading. Any lengthy response should probably not be colored with a red background.
func errorLabel(message string) {
	ct.ChangeColor(ct.White, true, ct.Red, false)
	fmt.Print(message)
	ct.ResetColor()
	fmt.Println("")
}

// Changes the color for the messages to green for success.
func successLabel(message string) {
	ct.Foreground(ct.Green, true)
	fmt.Print(message)
	ct.ResetColor()
	fmt.Println("")
}

// GetDiscfgNameFromFile simply returns the name of the set discfg name (TODO: will need to change as .discfg gets more complex).
func GetDiscfgNameFromFile() string {
	name := ""
	currentCfg, err := ioutil.ReadFile(DiscfgFileName)
	if err == nil {
		name = string(currentCfg)
	}
	return name
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

// Checks and formats the key name.
func formatKeyName(key string) (string, error) {
	var err error
	k := ""
	if len(key) > 0 {
		k = key
	} else {
		return "", errors.New(MissingKeyNameMsg)
	}

	// Ensure valid characters
	r, _ := regexp.Compile(`[\w\/\-]+$`)
	if !r.MatchString(k) {
		return "", errors.New(InvalidKeyNameMsg)
	}

	// Remove any trailing slashes (unless there's only one, the root).
	// NOTE: A tree structure is not yet supported. The user can define one, but there are no recursive features when getting/deleting.
	// This may come in a future version, for now the structure is flat. However, convention set by other tools (along with REST API endpoints)
	// makes using slashes a natural fit and discfg will assume they are being used. It could be thought of as a namespace.
	if len(k) > 1 {
		for k[len(k)-1:] == "/" {
			k = k[:len(k)-1]
		}
	}

	return k, err
}

func isJSONString(s string) bool {
	var js string
	err := json.Unmarshal([]byte(s), &js)
	return err == nil
}
func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

// FormatJSONValue sets the Item Value (an interface{}) as a map[string]interface{} (from []byte, which is how it's stored)
// so that it can be converted to JSON in an HTTP response. If it can't be represented in a map, then it'll be set as a string.
// For example, a string was set as the value, we can still represent that in JSON. However, if an image was stored...Then it's
// going to look ugly. It won't be a base64 string, it'll be the string representation of the binary data. Which apparently Chrome
// will render if given...But still. Not so hot. The user should know what they are setting and getting though and this should
// still technically return JSON with a usable value. Valid JSON at that. Just with some funny looking characters =)
func FormatJSONValue(resp config.ResponseObject) config.ResponseObject {
	// Don't attempt to Unmarshal or anything if the Value is empty. We wouldn't want to create a panic now.
	if resp.Item.Value == nil {
		return resp
	}
	resp.Item.Value = string(resp.Item.Value.([]byte))

	// The value could be base64 encoded, but it need not be.
	val, err := base64.StdEncoding.DecodeString(resp.Item.Value.(string)) //`eyJ1cGRhdGVkIjogImZyaWRheSJ9`)
	if err == nil && val != nil {
		resp.Item.Value = string(val)
	}

	var jsonData map[string]interface{}
	// Back to byte (because it could have potentially been base64 encoded)
	err = json.Unmarshal([]byte(resp.Item.Value.(string)), &jsonData)
	if err == nil {
		resp.Item.Value = jsonData
	}

	return resp
}
