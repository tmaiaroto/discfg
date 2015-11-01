package database

import (
	//"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sdming/gosnow"
	"github.com/tmaiaroto/discfg/config"
	//"log"
	"strconv"
)

// Each shipper has a struct which implements the Shipper interface.
type DynamoDB struct {
}

// Generates a new Snowflake
func generateId() int64 {
	// TODO: return the error so we can do something. Maybe this function isn't even needed...
	v, _ := gosnow.Default()

	// Alternatively you can set the worker id if you are running multiple snowflakes
	// TODO
	// v, err := gosnow.NewSnowFlake(100)

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
			{ // Required
				AttributeName: aws.String("k"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
			// More values...
		},
		KeySchema: []*dynamodb.KeySchemaElement{ // Required
			// One is required, but we use both a HASH (key name) and a RANGE (Snowflake).
			{
				AttributeName: aws.String("k"),    // Required
				KeyType:       aws.String("HASH"), // Required
			},
			{
				AttributeName: aws.String("id"),    // Required
				KeyType:       aws.String("RANGE"), // Required
			},
		},
		// GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		// 	{
		// 		IndexName: aws.String("KeyIndex"),
		// 		KeySchema: []*dynamodb.KeySchemaElement{
		// 			{
		// 				AttributeName: aws.String("k"),
		// 				KeyType:       aws.String("HASH"),
		// 			},
		// 		},
		// 		Projection: &dynamodb.Projection{
		// 			NonKeyAttributes: []*string{
		// 				aws.String("k"),
		// 			},
		// 			ProjectionType: aws.String("INCLUDE"),
		// 		},
		// 		// Annoying we need to provision for this too...
		// 		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
		// 			ReadCapacityUnits:  aws.Int64(1),
		// 			WriteCapacityUnits: aws.Int64(1),
		// 		},
		// 	},
		// },

		// Hard to estimate really. Should be passed along via command line when creating a new config.
		// Along with the table name. This will let people choose. Though it's kinda annoying someone must
		// think about this...
		// http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ProvisionedThroughputIntro.html
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ // Required
			ReadCapacityUnits:  aws.Int64(1), // Required
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

	// log.Println("Setting on table name: " + name)
	// log.Println(key)
	// log.Println(value)

	sid := generateId()
	sidStr := strconv.FormatInt(sid, 10)

	params := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
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
			//
			// keys need not be defined on CreateTable() - technically only the primary key need be defined upon table creation.
			"id": {
				N: aws.String(sidStr),
			},
			"k": { // Required
				S: aws.String(key),
			},
			// String for now...I mean everything is coming in as a string through stdin. Which makes JSON attractive. Though I'd like to get a little fancier...
			// Things like dates. Though MongoDB manages just fine with that...And discfg is only querying by key. So there's no much need...But who knows.
			// I really dig that we can store binary data into DynamoDB too. That might be special.
			"v": {
				S: aws.String(value),
			},
		},
		TableName: aws.String(name), // Required
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
	response, err := svc.PutItem(params)
	if err == nil {
		success = true
	}

	return success, response.String(), err
}

// Query results are always sorted by the range attribute value. If the data type of the range attribute is Number, the results are returned in numeric order;
// otherwise, the results are returned in order of ASCII character code values. By default, the sort order is ascending. To reverse the order, set the
// ScanIndexForward parameter to false.

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
		// Because "key" (and "value" for that matter) are reserved words
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":keyname": { // Required
				S: aws.String(key),
			},
		},
		KeyConditionExpression: aws.String("k = :keyname"),
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
			result["id"] = *response.Items[0]["id"].N
			result["value"] = *response.Items[0]["v"].S
		}
	}

	return success, result, err
}
