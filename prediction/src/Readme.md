# Prediction Model

Connection to producer and predict the stores, than save to data to influx.

## Local development

- install docker from https://docs.docker.com/get-docker/

- go to folder devops/dev

- update docker-compose.yaml line 38 with your local path to ssh keys

- run docker-compose up

- create ~/.netrc with content

machine machine gitlab.eaineu.com
login ales@storepredictor.com
password glpat-K7tZ-sq9Zxk9Due9d142

- set env GOPRIVATE=gitlab.eaineu.com/storepredictor

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
