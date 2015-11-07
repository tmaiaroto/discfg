package database

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tmaiaroto/discfg/config"
	"strconv"
	"strings"
	"time"
	// "log"
)

// Each shipper has a struct which implements the Shipper interface.
type DynamoDB struct {
}

// Configures DynamoDB service to use
func Svc(cfg config.Config) *dynamodb.DynamoDB {
	awsConfig := &aws.Config{Region: aws.String(cfg.Storage.DynamoDB.Region)}
	// Look in a variety of places for AWS credentials. First, try the credentials file set by AWS CLI tool.
	// Note the empty string instructs to look under default file path (different based on OS).
	// This file can have multiple profiles and a default profile will be used unless otherwise configured.
	// See: https://godoc.org/github.com/aws/aws-sdk-go/aws/credentials#SharedCredentialsProvider
	creds := credentials.NewSharedCredentials("", cfg.Storage.DynamoDB.CredProfile)
	_, err := creds.Get()
	// If that failed, try environment variables.
	if err != nil {
		// The following are checked:
		// Access Key ID: AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
		// Secret Access Key: AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
		creds = credentials.NewEnvCredentials()
	}

	// If credentials were passed via config, then use those. They will take priority over other methods.
	if cfg.Storage.DynamoDB.AccessKeyId != "" && cfg.Storage.DynamoDB.SecretAccessKey != "" {
		creds = credentials.NewStaticCredentials(cfg.Storage.DynamoDB.AccessKeyId, cfg.Storage.DynamoDB.SecretAccessKey, "")
	}
	awsConfig.Credentials = creds

	return dynamodb.New(awsConfig)
}

// Creates a new table for a configuration
func (db DynamoDB) CreateConfig(cfg config.Config, tableName string) (bool, interface{}, error) {
	svc := Svc(cfg)
	success := false
	wu := cfg.Storage.DynamoDB.WriteCapacityUnits
	ru := cfg.Storage.DynamoDB.ReadCapacityUnits
	// Must be at least 1
	if wu < 1 {
		wu = int64(1)
	}
	// Also must be at least 1, default to 2
	if ru < 1 {
		ru = int64(2)
	}

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
			ReadCapacityUnits:  aws.Int64(ru), // Required
			WriteCapacityUnits: aws.Int64(wu), // Required
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
func (db DynamoDB) Update(cfg config.Config, name string, key string, value string) (bool, config.Node, error) {
	var err error
	svc := Svc(cfg)
	success := false
	result := config.Node{Key: key}

	// log.Println("Setting on table name: " + name)
	// log.Println(key)
	// log.Println(value)

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
	//log.Println(value)

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
		// KEY and VALUE are reserved words so the query needs to dereference them
		ExpressionAttributeNames: map[string]*string{
			//"#k": aws.String("key"),
			"#v": aws.String("value"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			// value (always a string)
			":value": {
				B: []byte(value),
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

	// Conditional write operation (CAS)
	if cfg.ConditionalValue != "" {
		params.ExpressionAttributeValues[":condition"] = &dynamodb.AttributeValue{B: []byte(cfg.ConditionalValue)}
		params.ConditionExpression = aws.String("#v = :condition")
	}

	response, err := svc.UpdateItem(params)
	if err == nil {
		success = true
	}

	// The old values
	if val, ok := response.Attributes["value"]; ok {
		result.Value = val.B
		result.Version, _ = strconv.ParseInt(*response.Attributes["version"].N, 10, 64)
	}

	return success, result, err
}

// Gets a key in DynamoDB
func (db DynamoDB) Get(cfg config.Config, name string, key string) (bool, config.Node, error) {
	var err error
	svc := Svc(cfg)
	success := false
	result := config.Node{Key: key}

	params := &dynamodb.QueryInput{
		TableName: aws.String(name),

		// KEY and VALUE are reserved words so the query needs to dereference them
		ExpressionAttributeNames: map[string]*string{
			"#k": aws.String("key"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":key": {
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
		return success, result, err
	} else {
		success = true
		if len(response.Items) > 0 {
			// Every field should now be checked because it's possible to have a response without a value or verison.
			// For example, the root key "/" may only hold information about the config version and modified time.
			// It may not have a set value and therefore it also won't have a relative version either.
			// TODO: Maybe it should? We can always version it as 1 even if empty value. Perhaps also an empty string value...
			// But the update config version would need to have a compare for an empty value. See if DynamoDB can do that.
			// For now, just check the existence of keys in the map.
			if val, ok := response.Items[0]["value"]; ok {
				result.Value = val.B
			}
			if val, ok := response.Items[0]["version"]; ok {
				result.Version, _ = strconv.ParseInt(*val.N, 10, 64)
			}

			// If cfgVersion and cfgModified are set because it's the root key "/" then set those too.
			// This is only returned for the root key. no sense in making a separate get function because operations like
			// exporting would then require more queries than necessary. However, it won't be displayed in the node's JSON output.
			if val, ok := response.Items[0]["cfgVersion"]; ok {
				result.CfgVersion, _ = strconv.ParseInt(*val.N, 10, 64)
			}
			if val, ok := response.Items[0]["cfgModified"]; ok {
				result.CfgModifiedNanoseconds, _ = strconv.ParseInt(*val.N, 10, 64)
			}
		}
	}

	return success, result, err
}

// Deletes a key in DynamoDB
func (db DynamoDB) Delete(cfg config.Config, name string, key string) (bool, config.Node, error) {
	var err error
	svc := Svc(cfg)
	success := false
	result := config.Node{Key: key}

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
		TableName:    aws.String(name),
		ReturnValues: aws.String("ALL_OLD"),
		// TODO: think about this for statistics
		// INDEXES | TOTAL | NONE
		//ReturnConsumedCapacity: aws.String("ReturnConsumedCapacity"),
	}

	// Conditional delete operation
	if cfg.ConditionalValue != "" {
		// Alias value since it's a reserved word
		params.ExpressionAttributeNames = make(map[string]*string)
		params.ExpressionAttributeNames["#v"] = aws.String("value")
		// Set the condition expression value and compare
		params.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		params.ExpressionAttributeValues[":condition"] = &dynamodb.AttributeValue{B: []byte(cfg.ConditionalValue)}
		params.ConditionExpression = aws.String("#v = :condition")
	}

	response, err := svc.DeleteItem(params)
	if err != nil {
		return success, result, err
	} else {
		success = true
		if len(response.Attributes) > 0 {
			result.Value = response.Attributes["value"].B
			result.Version, _ = strconv.ParseInt(*response.Attributes["version"].N, 10, 64)
		}
	}

	return success, result, err
}

// Updates the configuration's global version and modified timestamp (fields unique to the root key "/")
func (db DynamoDB) UpdateConfigVersion(cfg config.Config, name string) bool {
	svc := Svc(cfg)
	success := false
	now := time.Now()
	params := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String("/"),
			},
		},
		TableName: aws.String(name),
		ExpressionAttributeNames: map[string]*string{
			"#m": aws.String("cfgModified"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			// modified timestamp (DynamoDB has no date type)
			":modified": {
				N: aws.String(strconv.FormatInt(now.UnixNano(), 10)),
			},
			// version increment
			":i": {
				N: aws.String("1"),
			},
		},
		ReturnValues:     aws.String("NONE"),
		UpdateExpression: aws.String("SET #m = :modified ADD cfgVersion :i"),
	}
	_, err := svc.UpdateItem(params)
	if err == nil {
		success = true
	}
	return success
}
