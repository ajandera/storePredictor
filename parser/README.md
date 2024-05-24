# Parser

Parser for open data.

## Local development

- install docker from https://docs.docker.com/get-docker/

- go to folder docker

- run docker-compose up

## Documentation

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
