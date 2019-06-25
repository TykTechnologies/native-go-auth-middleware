# go-dynamodb-basicauth-plugin

Built to be run natively as a package by Tyk Gateways.  

This is coded is intended to run as the "custom auth" part of the request lifecycle.

This will authenticate requests by connecting to DynamoDB and checking the Basic Auth credentials in the request to see if they match what's in the DB.  If success, will let the request continue, otherwise will return an auth error.

# Setup AWS credentials in the GO Plugin
In here, 
```
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", ""),
	})
```

Replace AKID and SECRET with your AWS credentials. Leave the third parameter exactly as is, an empty string.  Make sure the region is accurate also.

# Generate the Binary file
In the root of the "main.go" file, run 

`go build -o ./middleware/go/main.so -buildmode=plugin ./middleware/go`

Put the generated file somewhere Tyk Gateway can access it

# Build Tyk to run GOPLUGINS
`go build -tags 'goplugin' -o tyk .`

Then run the compiled Tyk

# Setup your API
in API Designer, click on "Raw API Definition"
1. Set ` "driver": "goplugin"` and `use_go_plugin_auth": true`
2. Set the request lifecycle to be run
```
"custom_middleware": {
      "auth_check": [
        {
          "name": "DynamoDBAuth",
          "path": "./middleware/go/main.so"
        }
      ],
   ```   
"name" has to be the name of the GO function

"path" is wherever you put the binary generated in step 1

# Testing

1. Create an entry in DynamoDB
```
username - hash
-------------
foo - bar
```

2. Base encode our user:password combo into "user:pwd" convention

`foo:bar` becomes `Zm9vOmJhcg==`

3. Make cURL to Tyk Gateway.
```curl http://www.tyk-test.com:8080/custom-auth-goplugin/headers -H "Authorization:Zm9vOmJhcg=="```

4.  Watch it rate limit!
