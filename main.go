package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	//"log"
	"time"
	//"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/spf13/cobra"
	//"github.com/tmaiaroto/discfg/storage"
	"github.com/tmaiaroto/discfg/config"
)

var _ time.Duration
var _ bytes.Buffer

var Config = config.Config{StorageInterfaceName: "dynamodb", Version: "0.1.0"}

var DiscfgCmd = &cobra.Command{
	Use:   "discfg",
	Short: "discfg is a distributed configuration service",
	Long:  `A distributed configuration service`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "discfg version number",
	Long:  `Displays the version number for discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("discfg v" + Config.Version)
	},
}

// TODO: Break out the functions for Run: ... but don't worry about organizing in such a manner that an AWS Lambda could call these directly.
// Lambda will get the same discfg binary that users use on their local machines...That's actually pretty nice for consistency in itself...
// But the bigger reason is less to maintain. No special provisions for Lambda.
// AWS Lambda can simply use the Node.js wrapper to make the exec call with command line options.
// It has to call the binary anyway...So why not just pass the necessary options to create new tables, add config keys, etc?
// Even if AWS Lambda supports Go in the future...Simply make something that proxies params from an API call to CLI options.
// Very easy. Pretend Lambda is just some person using the CLI tool.
// We aren't relying on AWS CLI... That's what the Go SDK is for. So all we need is config values...Which by default come from ~/.aws/
// Which won't exist on Lambda, but that's ok. Allow them to be passed in. Or come from a JSON file.
// Yes, a configuration for discfg, but it need not be distributed =) Hell, it could even come from a param passed to the API.
// So a Lambda could serve multiple accounts. Though SSL or no SSL, I might not pass my AWS IAM credentials through an HTTP RESTful interface.
//
// There will even be commands to setup Lambda and API Gateway...Those simply won't work once on AWS. Or maybe they will? Well, they need
// to be connected to an API Gateway resource to be called with the proper CLI parameters. So no, they just won't get used once on AWS.
//
// ... Still. Pull out the functions so they can more easily be tested and/or benchmarked.

// Highly recommended to only use `use` and `which` from CLI as a human.
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use a specific discfg",
	Long:  `For the current path, always use a specific discfg`,
	Run:   use,
}
var whichCmd = &cobra.Command{
	Use:   "which",
	Short: "shows current discfg in use",
	Long:  `Shows which discfg is currently selected for use at the current path`,
	Run:   which,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create config",
	Long:  `Creates a new discfg distributed configuration`,
	Run:   createCfg,
}
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set key value",
	Long:  `Sets a key value for a given discfg`,
	Run:   setKey,
}

func main() {
	// Check config

	// Set up commands
	DiscfgCmd.AddCommand(versionCmd)
	DiscfgCmd.PersistentFlags().StringVarP(&Config.OutputFormat, "format", "f", "human", "Output format for responses (human|json|slient)")

	DiscfgCmd.PersistentFlags().StringVarP(&Config.Storage.Region, "region", "r", "us-east-1", "AWS Region to use")
	//DiscfgCmd.PersistentFlags().StringVarP(&Config.CfgName, "cfgName", "n", "cfg", "The configuration name (cfg)")
	DiscfgCmd.AddCommand(useCmd)
	DiscfgCmd.AddCommand(whichCmd)
	DiscfgCmd.AddCommand(createCmd)
	DiscfgCmd.AddCommand(setCmd)
	DiscfgCmd.Execute()

	//storeField(svc)
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

// Ultimately this will be a lot less code. The examples are showing all of the possible options.
func storeField(svc *dynamodb.DynamoDB) {
	params := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{ // Required
			"field1": { // Required
				// B:    []byte("PAYLOAD"),
				// BOOL: aws.Bool(true),
				// BS: [][]byte{
				// 	[]byte("PAYLOAD"), // Required
				// 	// More values...
				// },
				// L: []*dynamodb.AttributeValue{
				// 	{ // Required
				// 	// Recursive values...
				// 	},
				// 	// More values...
				// },
				// M: map[string]*dynamodb.AttributeValue{
				// 	"Key": { // Required
				// 	// Recursive values...
				// 	},
				// 	// More values...
				// },
				// N: aws.String("NumberAttributeValue"),
				// NS: []*string{
				// 	aws.String("NumberAttributeValue"), // Required
				// 	// More values...
				// },
				// NULL: aws.Bool(true),
				// These are all the types...Here's the one we need in this case. A string.
				S: aws.String("test another"),
				// SS: []*string{
				// 	aws.String("StringAttributeValue"), // Required
				// 	// More values...
				// },
			},
			// keys need not be defined on CreateTable() - technically only the primary key need be defined upon table creation.
			// so, it's schemaless. which is perfect.
			"field2": {
				S: aws.String("yet more data"),
			},
			// More values...
		},
		TableName: aws.String("TestTableNameShift8"), // Required
		// Optional stuff...Will want to use the condtional operations though...
		//
		// ConditionExpression: aws.String("ConditionExpression"),
		// ConditionalOperator: aws.String("ConditionalOperator"),
		// Expected: map[string]*dynamodb.ExpectedAttributeValue{
		// 	"Key": { // Required
		// 		AttributeValueList: []*dynamodb.AttributeValue{
		// 			{ // Required
		// 				B:    []byte("PAYLOAD"),
		// 				BOOL: aws.Bool(true),
		// 				BS: [][]byte{
		// 					[]byte("PAYLOAD"), // Required
		// 					// More values...
		// 				},
		// 				L: []*dynamodb.AttributeValue{
		// 					{ // Required
		// 					// Recursive values...
		// 					},
		// 					// More values...
		// 				},
		// 				M: map[string]*dynamodb.AttributeValue{
		// 					"Key": { // Required
		// 					// Recursive values...
		// 					},
		// 					// More values...
		// 				},
		// 				N: aws.String("NumberAttributeValue"),
		// 				NS: []*string{
		// 					aws.String("NumberAttributeValue"), // Required
		// 					// More values...
		// 				},
		// 				NULL: aws.Bool(true),
		// 				S:    aws.String("StringAttributeValue"),
		// 				SS: []*string{
		// 					aws.String("StringAttributeValue"), // Required
		// 					// More values...
		// 				},
		// 			},
		// 			// More values...
		// 		},
		// 		ComparisonOperator: aws.String("ComparisonOperator"),
		// 		Exists:             aws.Bool(true),
		// 		Value: &dynamodb.AttributeValue{
		// 			B:    []byte("PAYLOAD"),
		// 			BOOL: aws.Bool(true),
		// 			BS: [][]byte{
		// 				[]byte("PAYLOAD"), // Required
		// 				// More values...
		// 			},
		// 			L: []*dynamodb.AttributeValue{
		// 				{ // Required
		// 				// Recursive values...
		// 				},
		// 				// More values...
		// 			},
		// 			M: map[string]*dynamodb.AttributeValue{
		// 				"Key": { // Required
		// 				// Recursive values...
		// 				},
		// 				// More values...
		// 			},
		// 			N: aws.String("NumberAttributeValue"),
		// 			NS: []*string{
		// 				aws.String("NumberAttributeValue"), // Required
		// 				// More values...
		// 			},
		// 			NULL: aws.Bool(true),
		// 			S:    aws.String("StringAttributeValue"),
		// 			SS: []*string{
		// 				aws.String("StringAttributeValue"), // Required
		// 				// More values...
		// 			},
		// 		},
		// 	},
		// 	// More values...
		// },
		// ExpressionAttributeNames: map[string]*string{
		// 	"Key": aws.String("AttributeName"), // Required
		// 	// More values...
		// },
		// ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
		// 	"Key": { // Required
		// 		B:    []byte("PAYLOAD"),
		// 		BOOL: aws.Bool(true),
		// 		BS: [][]byte{
		// 			[]byte("PAYLOAD"), // Required
		// 			// More values...
		// 		},
		// 		L: []*dynamodb.AttributeValue{
		// 			{ // Required
		// 			// Recursive values...
		// 			},
		// 			// More values...
		// 		},
		// 		M: map[string]*dynamodb.AttributeValue{
		// 			"Key": { // Required
		// 			// Recursive values...
		// 			},
		// 			// More values...
		// 		},
		// 		N: aws.String("NumberAttributeValue"),
		// 		NS: []*string{
		// 			aws.String("NumberAttributeValue"), // Required
		// 			// More values...
		// 		},
		// 		NULL: aws.Bool(true),
		// 		S:    aws.String("StringAttributeValue"),
		// 		SS: []*string{
		// 			aws.String("StringAttributeValue"), // Required
		// 			// More values...
		// 		},
		// 	},
		// 	// More values...
		// },
		//
		// The following will return info...
		// Needs to be one of:  [INDEXES, TOTAL, NONE]
		// ReturnConsumedCapacity:      aws.String("ReturnConsumedCapacity"),

		// Needs to be one of: [SIZE, NONE]
		// ReturnItemCollectionMetrics: aws.String("ReturnItemCollectionMetrics"),

		// Needs to be one of the following strings: [ALL_NEW, UPDATED_OLD, ALL_OLD, NONE, UPDATED_NEW]
		// ALL_OLD: Will return the previous value...because this PutItem will overwrite existing values.
		ReturnValues: aws.String("ALL_OLD"),
	}
	resp, err := svc.PutItem(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
