package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {}

func DynamoDBAuth(w http.ResponseWriter, r *http.Request) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", ""),
	})
	if err != nil {
		fmt.Println("Couldn't get AWS access")
		returnNoAuth(w, "Internal Error")
		return
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	if svc == nil {
		fmt.Println("Couldn't create new DynamoDB session")
		returnNoAuth(w, "Internal Error")
		return
	}

	username, password := unmarshalBasicAuth(r.Header.Get("Authorization"))

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
		returnNoAuth(w, "Internal Error")
		return
	} else if result.Item == nil {
		returnNoAuth(w, "Username not found.")
		return
	}

	basicAuth := BasicAuth{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &basicAuth)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	// Check Password
	if password != basicAuth.Hash {
		returnNoAuth(w, "Wrong Password.")
		return
	}

	fmt.Println(result)
}

func returnNoAuth(w http.ResponseWriter, errorMessage string) {
	jsonData, err := json.Marshal(errorMessage)
	if err != nil {
		fmt.Println("Couldn't marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func unmarshalBasicAuth(s string) (string, string) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Println("decode error:", err)
		panic("not decodable")
	}
	splitStr := strings.Split(string(decoded), ":")
	return string(splitStr[0]), string(splitStr[1])
}

type BasicAuth struct {
	Username string
	Hash     string
}
