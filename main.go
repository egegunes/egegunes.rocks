package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
)

type comment struct {
	Name      string
	Message   string
	CreatedAt string
}

func (c comment) DisplayDate() string {
	createdAt, err := time.Parse(time.RFC3339, c.CreatedAt)
	if err != nil {
		return c.CreatedAt
	}
	return createdAt.Format("2006-01-02 15:04")
}

func createTable(db *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("CreatedAt"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("CreatedAt"),
				KeyType:       aws.String("RANGE"),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := db.CreateTable(input)

	if err != nil {
		return err
	}

	return nil
}

func tableExists(db *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	_, err := db.DescribeTable(input)

	return err
}

func getComments(db *dynamodb.DynamoDB, tableName string) ([]comment, error) {
	var comments []comment

	input := &dynamodb.ScanInput{
		ConsistentRead: aws.Bool(true),
		TableName:      aws.String(tableName),
	}
	result, err := db.Scan(input)

	if err != nil {
		return comments, err
	}

	dynamodbattribute.UnmarshalListOfMaps(result.Items, &comments)

	return comments, err
}

func putComment(db *dynamodb.DynamoDB, tableName string, c comment) error {
	av, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)

	return err
}

func main() {
	db := dynamodb.New(session.New(), &aws.Config{
		Region:   aws.String("eu-central-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	tableName := "Comments"

	if err := tableExists(db, tableName); err != nil {
		log.Println("table does not exist, creating new")
		err := createTable(db, tableName)
		if err != nil {
			log.Printf("can't create table: %s", err)
		}
	}

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.tmpl")
	router.GET("/", func(c *gin.Context) {
		comments, err := getComments(db, tableName)
		if err != nil {
			log.Printf("can't get comments: %s", err)
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"comments": comments})
	})
	router.POST("/", func(c *gin.Context) {
		name := c.PostForm("name")
		message := c.PostForm("message")
		err := putComment(db, tableName, comment{name, message, time.Now().Format(time.RFC3339)})
		if err != nil {
			log.Printf("can't put comment: %s", err)
		}
		comments, err := getComments(db, tableName)
		if err != nil {
			log.Printf("can't get comments: %s", err)
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"comments": comments})
	})
	router.Run(":8080")
}
