package dbi

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	uuid "github.com/nu7hatch/gouuid"
	perrors "github.com/pkg/errors"
)

//Client minimal dynamodb client
type Client interface {
	CreateTable(input *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	ScanPages(input *dynamodb.ScanInput, fn func(*dynamodb.ScanOutput, bool) bool) error
}

//Database standard Interface implementation
type Database struct {
	Client Client
}

func (db *Database) CreateTable(name string) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("uuid"),
				KeyType:       aws.String("HASH"),
			},
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("moment"),
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("uuid"),
				AttributeType: aws.String("S"),
			},
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("moment"),
				AttributeType: aws.String("S"),
			},
		},
	}

	_, err := db.Client.CreateTable(input)
	return perrors.Wrapf(err, "failed creating table `%s`", name)
}

func (db *Database) Index(tableName string) (map[string]string, error) {
	// var uuid, data string
	// result := map[string]string{}

	// rows, err := db.DB.Query("SELECT (uuid, data) FROM ?", tableName)
	// if err != nil {
	// 	return nil, perrors.Wrapf(err, "failed to index table `%s`", tableName)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	err = rows.Scan(&uuid, &data)
	// 	if err != nil {
	// 		return nil, perrors.Wrapf(err, "failed to scan row of table `%s`", tableName)
	// 	}

	// 	result[uuid] = data
	// }

	// return result, nil
	return nil, nil
}

func (db *Database) Insert(tableName string, document map[string]string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return perrors.Wrapf(err, "failed generating UUIDv4 for table `%s`", tableName)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"uuid": &dynamodb.AttributeValue{
				S: aws.String(uuid.String()),
			},
			"moment": &dynamodb.AttributeValue{
				S: aws.String(time.Now().Format(time.RFC3339)),
			},
		},
	}

	for k, v := range document {
		input.Item[k] = &dynamodb.AttributeValue{S: aws.String(v)}
	}

	_, err = db.Client.PutItem(input)
	return perrors.Wrapf(err, "failed inserting into `%s`", tableName)
}
