package main

import (
	"encoding/json"
	"github.com/apex/go-apex"
	"github.com/tmaiaroto/discfg/commands"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/version"
	"os"
	"time"
)

// To change these settings for DynamoDB, deploy with a different environment variable.
// apex deploy -s DISCFG_DB_REGION=us-west-1
var discfgDBRegion = os.Getenv("DISCFG_REGION")
var discfgDBTable = os.Getenv("DISCFG_TABLE")

// The JSON message passd to the Lambda (should include key, value, etc.)
type message struct {
	Name string `json:"name"`
	// Comes in as string, but needs to be converted to int64
	TTL   string `json:"ttl"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Raw   string `json:"raw"`
}

var options = config.Options{StorageInterfaceName: "dynamodb", Version: version.Semantic}

func main() {
	// If not set for some reason, use us-east-1 by default.
	if discfgDBRegion == "" {
		discfgDBRegion = "us-east-1"
	}

	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		options.Storage.AWS.Region = discfgDBRegion
		// Each discfg API can be configured with a default table name.
		options.CfgName = discfgDBTable
		// Overwritten by the message passed to the Lambda.
		if m.Name != "" {
			options.CfgName = m.Name
		}
		options.Key = m.Key

		resp := commands.GetKey(options)

		// Format the expiration time (if applicable). This prevents output like "0001-01-01T00:00:00Z" when empty
		// and allows for the time.RFC3339Nano format to be used whereas time.Time normally marshals to a different format.
		if resp.Item.TTL > 0 {
			resp.Item.OutputExpiration = resp.Item.Expiration.Format(time.RFC3339Nano)
		}

		r := commands.FormatJSONValue(resp)
		return r, nil
	})
}
