// Build with: GOOS=linux GOARCH=amd64 go build main.go
package main

import (
	"encoding/json"
	"github.com/jasonmoo/lambda_proc"
	"github.com/tmaiaroto/discfg/commands"
	"github.com/tmaiaroto/discfg/config"
	"strconv"
)

var Config = config.Options{StorageInterfaceName: "dynamodb", Version: "0.2.0"}

func main() {
	lambda_proc.Run(func(context *lambda_proc.Context, eventJSON json.RawMessage) (interface{}, error) {
		var v map[string]interface{}
		if err := json.Unmarshal(eventJSON, &v); err != nil {
			return nil, err
		}
		var out interface{}
		out = nil

		// Set us-east-1 as the default region.
		Config.Storage.DynamoDB.Region = "us-east-1"
		Config.Storage.DynamoDB.ReadCapacityUnits = 1
		Config.Storage.DynamoDB.WriteCapacityUnits = 2
		creds := v["creds"].(map[string]interface{})
		if accessKeyId, ok := creds["accessKeyId"].(string); ok {
			Config.Storage.DynamoDB.AccessKeyId = accessKeyId
		}
		if secretAccessKey, ok := creds["secretAccessKey"].(string); ok {
			Config.Storage.DynamoDB.SecretAccessKey = secretAccessKey
		}
		if sessionToken, ok := creds["sessionToken"].(string); ok {
			Config.Storage.DynamoDB.SessionToken = sessionToken
		}

		// No format. Return the ResponeObject
		Config.OutputFormat = ""

		// Adjust Config based on event JSON.
		if region, ok := v["region"].(string); ok {
			Config.Storage.DynamoDB.Region = region
		}

		if cfgName, ok := v["name"].(string); ok {
			Config.CfgName = cfgName
		}

		if key, ok := v["key"].(string); ok {
			Config.Key = key
		}

		if value, ok := v["value"].(string); ok {
			Config.Value = value
		}

		if ttl, ok := v["ttl"].(string); ok {
			ttlInt, err := strconv.ParseInt(ttl, 10, 64)
			if err == nil {
				Config.TTL = ttlInt
			}
		}

		if condVal, ok := v["condition"].(string); ok {
			Config.ConditionalValue = condVal
		}

		// Determine the command to run.
		if command, ok := v["command"].(string); ok {
			switch command {
			case "get":
				resp := commands.GetKey(Config)
				formattedResp := commands.Out(Config, resp)
				j, _ := json.Marshal(&formattedResp)
				out = string(j)
				//out = formattedResp
			case "set":
				resp := commands.SetKey(Config)
				formattedResp := commands.Out(Config, resp)
				j, _ := json.Marshal(&formattedResp)
				out = string(j)
				//out = formattedResp
			case "delete":
				resp := commands.DeleteKey(Config)
				formattedResp := commands.Out(Config, resp)
				j, _ := json.Marshal(&formattedResp)
				out = string(j)
				//out = formattedResp
			}
		}

		return out, nil
	})
}
