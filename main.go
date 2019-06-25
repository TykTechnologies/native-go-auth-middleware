package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/user"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc = &dynamodb.DynamoDB{}

// Run on startup by Tyk when loaded.  Bootstrapping the service here
func init() {
	// Authenticate User in AWS
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", ""),
	})
	if err != nil {
		fmt.Println("Couldn't get AWS access")
		panic("Internal Error")
	}
	// Create DynamoDB client
	svc = dynamodb.New(sess)
	if svc == nil {
		fmt.Println("Couldn't create new DynamoDB session")
		panic("Internal Error")
	}
}
func main() {}

// DynamoDBAuth - Main method to be run on each request
func DynamoDBAuth(w http.ResponseWriter, r *http.Request) {
	encodedHeaderValue := r.Header.Get("Authorization")
	username, password := unmarshalBasicAuth(encodedHeaderValue)

	// Get the Basic Auth user/pass matching the username in the request from DynamoDB
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
		returnNoAuth(w, "Internal Error")
		return
	}

	// Check Password
	if password != basicAuth.Hash {
		returnNoAuth(w, "Wrong Password.")
		return
	}

	// create session
	session := &user.SessionState{
		OrgID: "default",
		Rate:  5,
		Per:   10,
	}
	ctx.SetSession(r, session, encodedHeaderValue, true)

	// Let the request continue
	fmt.Println("Passed Auth")
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

// BasicAuth Looks same as the DynamoDB structure
type BasicAuth struct {
	Username string
	Hash     string
}
