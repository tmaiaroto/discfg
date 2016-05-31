package main

import (
	"encoding/json"
	"github.com/apex/go-apex"
	"github.com/tmaiaroto/discfg/commands"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/version"
	"os"
)

// To change these settings for DynamoDB, deploy with a different environment variable.
// apex deploy -s DISCFG_DB_REGION=us-west-1
var discfgDBRegion = os.Getenv("DISCFG_REGION")

// The JSON message passed to the Lambda
type message struct {
	Name string `json:"name"`
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

		if m.Name != "" {
			options.CfgName = m.Name
		}

		resp := commands.DeleteCfg(options)

		return commands.FormatJsonValue(resp), nil
	})
}
