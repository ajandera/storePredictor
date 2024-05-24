### Documentation

Run command from src folder

    godoc -http=:6060 

Install godoc if needed

    go install -v golang.org/x/tools/cmd/godoc@latest

### Tests

Run command from src folder
    
    go test -cover ./...

Generate coverage profile

    go test ./... -coverprofile coverage.out

See coverage by each function

    go tool cover -func coverage.out 

See coverage preview in html

    go tool cover -html coverage.out

### Apiary connector to generate API doc

https://help.apiary.io/tools/apiary-cli/

Install apiary cli

    gem install apiaryio

Apiary token generates from

https://login.apiary.io/tokens

Add to env 

    export APIARY_API_KEY="<token>"

Apiary API_NAME: storepredictorshopapi

Retrieve actual blueprint document

    apiary fetch --api-name="storepredictorshopapi"

Generate local files

    apiary fetch --api-name="<API_NAME>" --output="apiary.apib"
    # or
    apiary fetch --api-name="<API_NAME>" --output="swagger.yaml"

Preview locally

    apiary preview

Publish to Apiary

    apiary publish --api-name="storepredictorshopapi"

### Generate api documentation

Install SwagGo

    go get -u github.com/swaggo/swag/cmd/swag

Add GOPATH to your system path (mac os exmaple)

    export PATH=$(go env GOPATH)/bin:$PATH 

Go to you project and init swagGo, in root folder

    swag init --parseDependency

To update swagger only from some files

    swag init --parseDependency -g main.go

Locally will se on 

http://localhost:9999/swagger/index.html


Preview in Appiary

    apiary preview --path="/path/to/swagger.yaml"

Publish to Appiary use command from src/docs files

    apiary publish --api-name="storepredictorshopapi"

