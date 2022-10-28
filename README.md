# Welsh Academy
<a href="https://goreportcard.com/report/github.com/mjehanno/welsh-academy"><img src="https://goreportcard.com/badge/github.com/mjehanno/welsh-academy"/></a>


Welsh-Academy is a simple REST-API made to handle recipes and ingredients.

## Usage

### Build / Run

#### With Taskfile

You can use [Taskfile](https://taskfile.dev/) in order to "play" with the application. 
It will provide you some small command you can run to build, test, ... the application.

But you will first need to install this tool before being able to use the command that are in the file :

`go install github.com/go-task/task/v3/cmd/task@latest`

Action | command
--- | ---
Build| `task build`
Run| `task run`
Test | `task test`
Coverage | `task coverage`
Build the Docker image | `task docker-build`
Run the Docker container | `task docker-run`
Generate swagger documentation | `task gen-swagger`

#### Without Taskfile

You can do the same actions without using Task but commands are a bit longer

Action | command
--- | ---
Build| `go build -o ./bin/welsh-academy -v cmd/welsh-academy/main.go`
Run| `./bin/welsh-academy`
Test | `go test ./... -v -cover`
Coverage | `go test ./... -coverprofile=coverage.out`
Build the Docker image | `docker build . -t rest-document`
Run the Docker container | `docker run -p 9000:9000 rest-document`
Generate swagger documentation | `swag init -g cmd/welsh-academy/main.go --parseDependency --parseInternal`

### REST api

The api use the following HTTP Method :

Action | Method 
--- | ---
Retrieve one or many recipe/ingredient | GET
Create a recipe/ingredient/user | POST
Untag a favorite recipe | DELETE

Considered route are either

StatusCode :
- 200 => action did work
- 201 => object was created (POST request)
- 204 => content has been deleted
- 400 => error from user 
- 401 => need to login before
- 500 => error in the api


## Documentation

The api documentation is available as an OpenAPI spec file. 
This file is loaded as the application start and rendered on `:9000/docs/index.html`.


