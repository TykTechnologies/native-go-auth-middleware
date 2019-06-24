# go-dynamodb-basicauth-plugin

Built to be run natively as a package by Tyk Gateways.  

# Generate the Binary file
In the root of the "main.go" file, run 
`go build -o ./middleware/go/main.so -buildmode=plugin ./middleware/go`
Put the generated file somewhere Tyk Gateway can access it

# Build Tyk to run GOPLUGINS
`go build -tags 'coprocess grpc goplugin' -o tyk .`
you may only need goplugin in that list above

Then run the compiled Tyk

# Setup your API
in API Designer, click on "Raw API Definition"
1. Set ` "driver": "goplugin"`
2. Choose somewhere for your middleware to run in the cycle. ie:
"custom_middleware": {
      "pre": [
        {
          "name": "MyCustomPlugin",
          "path": "./middleware/go/helloworld.so"
        }
      ],
      
Pre is the phase in the cycle where it runs.
"name" has to be the name of the GO function
"path" is wherever you put the binary generated in step 1
