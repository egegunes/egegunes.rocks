package main

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

type comment struct {
	Message   string
	CreatedAt time.Time
}

func createTable(db *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Message"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Time"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Message"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Time"),
				KeyType:       aws.String("RANGE"),
			},
		},
		TableName: aws.String("Comments"),
	}
	_, err := db.CreateTable(input)

	if err != nil {
		return err
	}

	return nil
}

func tableExists(db *dynamodb.DynamoDB, tableName string) bool {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	_, err := db.DescribeTable(input)

	return err != nil
}

func main() {
	db := dynamodb.New(session.New(), &aws.Config{
		Region:   aws.String("eu-central-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	tableName := "Comments"

	if !tableExists(db, tableName) {
		createTable(db, tableName)
	}

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.tmpl")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"comments": []comment{
				comment{"deneme", time.Now()},
				comment{"deneme2", time.Now()},
			},
		})
	})
	router.Run(":8080")
}
