package database

import (
	//"errors"
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sdming/gosnow"
	"github.com/tmaiaroto/discfg/config"
	"log"
	//"strconv"
	//"encoding/json"
	"strings"
)

// Each shipper has a struct which implements the Shipper interface.
type DynamoDB struct {
}

// Generates a new Snowflake
func generateId() int64 {
	// TODO: return the error so we can do something. Maybe this function isn't even needed...
	// Or maybe return it in the format dynamodb wants using aws package N: { aws.String() } ...
	v, _ := gosnow.Default()
	id, _ := v.Next()

	return int64(id)
}

// Configures DynamoDB service to use
func Svc(cfg config.Config) *dynamodb.DynamoDB {
	return dynamodb.New(&aws.Config{Region: aws.String(cfg.Storage.Region)})
}

// Creates a new table for a configuration
func (db DynamoDB) CreateConfig(cfg config.Config, tableName string) (bool, interface{}, error) {
	svc := Svc(cfg)
	success := false

	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: aws.String("S"),
			},
			// {
			// 	AttributeName: aws.String("id"),
			// 	AttributeType: aws.String("N"),
			// },
			// // More values...
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			// One is required, but we use both a HASH (key name) and a RANGE (Snowflake).
			{
				AttributeName: aws.String("key"),
				KeyType:       aws.String("HASH"),
			},
			// {
			// 	AttributeName: aws.String("id"),
			// 	KeyType:       aws.String("RANGE"),
			// },
		},
		// Hard to estimate really. Should be passed along via command line when creating a new config.
		// Along with the table name. This will let people choose. Though it's kinda annoying someone must
		// think about this...
		// http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ProvisionedThroughputIntro.html
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ // Required
			ReadCapacityUnits:  aws.Int64(2), // Required
			WriteCapacityUnits: aws.Int64(1), // Required
		},
		TableName: aws.String(tableName), // Required

	}
	response, err := svc.CreateTable(params)
	if err == nil {
		tableStatus := *response.TableDescription.TableStatus
		if tableStatus == "CREATING" || tableStatus == "ACTIVE" {
			success = true
		}
	}
	return success, response, err
}

// Updates a key in DynamoDB
func (db DynamoDB) Update(cfg config.Config, name string, key string, value string) (bool, interface{}, error) {
	var err error
	svc := Svc(cfg)
	success := false
	result := make(map[string]string)

	// log.Println("Setting on table name: " + name)
	// log.Println(key)
	// log.Println(value)

	//sid := generateId()
	//sidStr := strconv.FormatInt(sid, 10)

	keys := strings.Split(key, "/")
	parents := []*string{}
	if len(keys) > 0 {
		// Keep appending previous path so each parent key is an absolute path
		prevKey := ""
		var buffer bytes.Buffer
		for i := range keys {
			// Don't take an empty value or itself as a parent
			if keys[i] != "" && keys[i] != key {
				buffer.WriteString(prevKey)
				buffer.WriteString("/")
				buffer.WriteString(keys[i])
				prevKey = buffer.String()
				parents = append(parents, aws.String(prevKey))
				buffer.Reset()
			}
		}
	}

	// TODO: Fix - the panic is when there are no child. parents slice has issues.
	// TODO: JSON seems to be saving...but check output formatting (unescaping, parsing - when possible)
	log.Println(value)

	// DynamoDB type cheat sheet:
	// B: []byte("some bytes")
	// BOOL: aws.Bool(true)
	// BS: [][]byte{[]byte("bytes and bytes")}
	// L: []*dynamodb.AttributeValue{{...recursive values...}}
	// M: map[string]*dynamodb.AttributeValue{"key": {...recursive...} }
	// N: aws.String("number")
	// NS: []*String{aws.String("number"), aws.String("number")}
	// NULL: aws.Bool(true)
	// S: aws.String("string")
	// SS: []*string{aws.String("string"), aws.String("string")}

	// If always putting new items, there's no conditional update.
	// But the only way to update is to make the items have a HASH only index instead of HASH + RANGE.

	params := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(name),
		// TODO: This. If passed from the CLI, this needs to be checked.
		// ConditionExpression: aws.String("ConditionExpression"),

		// KEY and VALUE are reserved words so the query needs to dereference them
		ExpressionAttributeNames: map[string]*string{
			//"#k": aws.String("key"),
			"#v": aws.String("value"),
		},

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			// value (always a string)
			":value": {
				S: aws.String(value),
				// Always store bytes?
				// B: []byte(value),
			},
			// parents
			":pv": {
				SS: parents,
			},
			// version increment
			":i": {
				N: aws.String("1"),
			},
		},

		//ReturnConsumedCapacity:      aws.String("TOTAL"),
		//ReturnItemCollectionMetrics: aws.String("ReturnItemCollectionMetrics"),
		ReturnValues:     aws.String("ALL_OLD"),
		UpdateExpression: aws.String("SET #v = :value, parents = :pv ADD version :i"),
	}
	response, err := svc.UpdateItem(params)

	// params := &dynamodb.PutItemInput{
	// 	Item: map[string]*dynamodb.AttributeValue{
	// 		"id": {
	// 			N: aws.String(sidStr),
	// 		},
	// 		"k": { // Required
	// 			S: aws.String(key),
	// 		},
	// 		// String for now...I mean everything is coming in as a string through stdin. Which makes JSON attractive. Though I'd like to get a little fancier...
	// 		// Things like dates. Though MongoDB manages just fine with that...And discfg is only querying by key. So there's no much need...But who knows.
	// 		// I really dig that we can store binary data into DynamoDB too. That might be special.
	// 		"v": {
	// 			S: aws.String(value),
	// 		},
	// 		// The parent path(s)
	// 		"parents": {
	// 			SS: parents,
	// 		},
	// 	},
	// 	TableName: aws.String(name), // Required
	// 	// The following will return info...
	// 	// Needs to be one of:  [INDEXES, TOTAL, NONE]
	// 	// ReturnConsumedCapacity:      aws.String("ReturnConsumedCapacity"),

	// 	// Needs to be one of: [SIZE, NONE]
	// 	// ReturnItemCollectionMetrics: aws.String("ReturnItemCollectionMetrics"),

	// 	// Needs to be one of the following strings: [ALL_NEW, UPDATED_OLD, ALL_OLD, NONE, UPDATED_NEW]
	// 	// ALL_OLD: Will return the previous value...because this PutItem will overwrite existing values.
	// 	ReturnValues: aws.String("ALL_OLD"),
	// }
	// response, err := svc.PutItem(params)
	if err == nil {
		success = true
	}

	// The old values
	if val, ok := response.Attributes["value"]; ok {
		result["value"] = *val.S
		result["version"] = *response.Attributes["version"].N
	}

	return success, result, err
}

// Gets a key in DynamoDB
func (db DynamoDB) Get(cfg config.Config, name string, key string) (bool, interface{}, error) {
	var err error
	svc := Svc(cfg)
	success := false
	result := make(map[string]string)

	// still not clear on what the difference between dynamodb.New() is and session.New()
	// may only matter if there were multiple queries...we aren't going to have "sessions" ... it's just one command - one query.
	//svc := dynamodb.New(session.New())

	params := &dynamodb.QueryInput{
		TableName: aws.String(name),

		// KEY and VALUE are reserved words so the query needs to dereference them
		ExpressionAttributeNames: map[string]*string{
			"#k": aws.String("key"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":key": { // Required
				S: aws.String(key),
			},
		},
		KeyConditionExpression: aws.String("#k = :key"),
		// TODO: Return more? It's nice to have a history now whereas previously I thought I might now have one...But what's the use?
		Limit: aws.Int64(1),

		// INDEXES | TOTAL | NONE (not required - not even sure if I need to worry about it)
		ReturnConsumedCapacity: aws.String("TOTAL"),
		// Important: This needs to be false so it returns results in descending order. If it's true (the default), it's sorted in the
		// order values were stored. So the first item stored for the key ever would be returned...But the latest item is needed.
		ScanIndexForward: aws.Bool(false),
		// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html#DDB-Query-request-Select
		Select: aws.String("ALL_ATTRIBUTES"),
	}
	response, err := svc.Query(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		//fmt.Println(err.Error())
		return success, nil, err
	} else {
		success = true
		if len(response.Items) > 0 {
			// result["id"] = *response.Items[0]["id"].N
			result["value"] = *response.Items[0]["value"].S
			result["version"] = *response.Items[0]["version"].N
		}
	}

	return success, result, err
}

// Deletes a key in DynamoDB
func (db DynamoDB) Delete(cfg config.Config, name string, key string) (bool, interface{}, error) {
	var err error
	svc := Svc(cfg)
	success := false
	result := make(map[string]string)

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{ // Required
			"key": {
				S: aws.String(key),
			},
		},
		TableName:    aws.String(name),
		ReturnValues: aws.String("ALL_OLD"),
		// TODO: allow an option to be passed for conditional delete (very nice feature to have)
		// ConditionExpression: aws.String("ConditionExpression"),
		//
		// TODO: think about this for statistics
		// INDEXES | TOTAL | NONE
		//ReturnConsumedCapacity: aws.String("ReturnConsumedCapacity"),
	}
	response, err := svc.DeleteItem(params)
	if err != nil {
		return success, nil, err
	} else {
		success = true
		if len(response.Attributes) > 0 {
			result["value"] = *response.Attributes["value"].S
			result["version"] = *response.Attributes["version"].N
		}
	}

	return success, result, err
}
