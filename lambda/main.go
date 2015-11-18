package main

import (
	"encoding/json"
	"flag"
	"github.com/jasonmoo/lambda_proc"
	"github.com/tmaiaroto/discfg/commands"
	"github.com/tmaiaroto/discfg/config"
	"log"
)

var Config = config.Options{StorageInterfaceName: "dynamodb", Version: "0.2.0"}

func main() {
	var testFlag = flag.Bool("test", false, "run a test command")
	flag.Parse()

	if *testFlag {
		test()
	} else {
		lambda_proc.Run(func(context *lambda_proc.Context, eventJSON json.RawMessage) (interface{}, error) {
			var v map[string]interface{}
			if err := json.Unmarshal(eventJSON, &v); err != nil {
				return nil, err
			}
			// Defaults...Change with values from message.
			Config.Storage.Region = "us-east-1"
			//return getLocation(v["source-ip"].(string))
			// switch on v["command"] or something...call the proper commands
			resp := commands.GetKey(Config, []string{"mycfg", "/num"})
			return resp, nil
		})
	}
}

// TODO: probably just remove this. need to find a way to test Lambdas locally.
func test() {
	Config.Storage.Region = "us-east-1"
	resp := commands.GetKey(Config, []string{"mycfg", "/num"})
	log.Println(string(resp.Node.Value))
}
