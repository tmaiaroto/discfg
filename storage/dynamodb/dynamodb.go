package database

import (
	//"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tmaiaroto/discfg/config"
	"log"
)

// Each shipper has a struct which implements the Shipper interface.
type DynamoDB struct {
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
				AttributeName: aws.String("key"),
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
				AttributeName: aws.String("key"),  // Required
				KeyType:       aws.String("HASH"), // Required
			},
			{
				AttributeName: aws.String("id"),    // Required
				KeyType:       aws.String("RANGE"), // Required
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("KeyIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("key"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					NonKeyAttributes: []*string{
						aws.String("key"),
					},
					ProjectionType: aws.String("INCLUDE"),
				},
				// Annoying we need to provision for this too...
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
		},

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
	//svc := Svc(cfg)
	success := false

	log.Println("Setting on table name: " + name)
	log.Println(key)
	log.Println(value)

	return success, nil, err
}
