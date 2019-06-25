package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	usersTable = "users"
	usernameField = "username"
	authorizationHeader = "Authorization"
	region = "eu-west-2"
)

// BasicAuth Looks same as the DynamoDB structure
type BasicAuth struct {
	Username string
	Hash     string
}

var svc *dynamodb.DynamoDB

// Run on startup.  Bootstrapping the service here
func init() {
	// Authenticate User in AWS
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
	})
	if err != nil {
		log.Fatalf("couldn't get AWS access: %v", err.Error())
	}
	// Create DynamoDB client
	svc = dynamodb.New(sess)
}

func main() {
	// uncomment and compile normally to test locally
	//http.HandleFunc("/", DynamoDBAuth)
	//log.Fatal(http.ListenAndServe(":8000", nil))
}

// Main method to be run by Tyk
func DynamoDBAuth(w http.ResponseWriter, r *http.Request) {
	username, password, err := unmarshalBasicAuth(r.Header.Get(authorizationHeader))
	if err != nil {
		returnNoAuth(w, err.Error())
		return
	}

	log.Printf("attempted access with %s:%s", username, password)

	// Get the Basic Auth user/pass matching the username in the request from DynamoDB
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(usersTable),
		Key: map[string]*dynamodb.AttributeValue{
			usernameField: {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		returnNoAuth(w, err.(awserr.Error).Message())
		return
	}

	basicAuth := BasicAuth{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &basicAuth)
	if err != nil {
		returnNoAuth(w, err.Error())
		return
	}

	if basicAuth.Username == "" {
		returnNoAuth(w, "User not found.")
		return
	}

	// Check Password
	if password != basicAuth.Hash {
		returnNoAuth(w, "Wrong Password.")
		return
	}

	// Let the request continue
	fmt.Println("Auth passed")
}

func returnNoAuth(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized) + " " + errorMessage))
}

func unmarshalBasicAuth(s string) (string, string, error) {
	if s == "" {
		return "", "", errors.New("no credentials supplied")
	}

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", "", err
		fmt.Println("decode error:", err)
	}

	splitStr := strings.Split(string(decoded), ":")
	if len(splitStr) != 2 {
		return "", "", errors.New("not in user:pass format")
	}

	return string(splitStr[0]), string(splitStr[1]), nil
}
