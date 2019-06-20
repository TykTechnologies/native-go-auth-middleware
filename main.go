package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {}

func DynamoDBAuth(w http.ResponseWriter, r *http.Request) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", ""),
	})
	if err != nil {
		fmt.Println("Couldn't get AWS access")
		return
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	if svc == nil {
		fmt.Println("Couldn't create new DynamoDB session")
		return
	}

	username := "foo" // read this from outside

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("basic-auth"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}

type BasicAuth struct {
	Username string
	Password string
}
