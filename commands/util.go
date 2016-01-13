// Various utilities used by commands are found in this file as well as response structs, constants, etc.
package commands

import (
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

// TODO: Change this and use config/status.go instead
const NotEnoughArgsMsg = "Not enough arguments passed. Run 'discfg help' for usage."
const ValueRequired = "A value is required. Run 'discfg help' for usage."
const DiscfgFileName = ".discfg"

// Output
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
	// if resp.Node.Value != nil {
	// 	if !isJSON(string(resp.Node.Value)) {
	// 		// Return base64 when not JSON?
	// 		// b64Str := base64.StdEncoding.EncodeToString(resp.Node.Value)
	// 		//resp.Node.Value = []byte(strconv.Quote(b64Str))
	// 		resp.Node.Value = []byte(strconv.Quote(string(resp.Node.Value)))
	// 	}
	// 	// The output value is always raw JSON. It is not stored in the data store.
	// 	// It's simply for display.
	// 	resp.Node.OutputValue = json.RawMessage(resp.Node.Value)
	// }

	// // Same for the PrevNode if set
	// if resp.PrevNode.Value != nil {
	// 	if !isJSON(string(resp.PrevNode.Value)) {
	// 		resp.PrevNode.Value = []byte(strconv.Quote(string(resp.PrevNode.Value)))
	// 	}
	// 	resp.PrevNode.OutputValue = json.RawMessage(resp.PrevNode.Value)
	// }

	// Format the expiration time (if applicable). This prevents output like "0001-01-01T00:00:00Z" when empty
	// and allows for the time.RFC3339Nano format to be used whereas time.Time normally marshals to a different format.
	if resp.Node.TTL > 0 {
		resp.Node.OutputExpiration = resp.Node.Expiration.Format(time.RFC3339Nano)
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
		if resp.Node.Value != nil {
			// The value should be a byte array, for th CLI we want a string.
			fmt.Println(string(resp.Node.Value.([]byte)))
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

// Just returns the name of the set discfg name (TODO: will need to change as .discfg gets more complex).
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
		return "", errors.New("Missing key name")
	}

	// Ensure valid characters
	r, _ := regexp.Compile(`[\w\/\-]+$`)
	if !r.MatchString(k) {
		return "", errors.New("Invalid key name")
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
