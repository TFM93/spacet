# SpaceT Microservice

## Table of Contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Usage](#usage)
5. [Environment Variables](#environment_variables)
6. [API](#api)
7. [Implementation Notes](#implementation)

## Introduction
SpaceT Microservice schedules launches by sharing the launchpads with spaceX.
I've considered that we can't control the spaceX launches therefore, the API must prioritize them. This means our internally booked launches can be canceled after booking if they conflict with SpaceX's schedule.
The system works by periodically fetching the upcoming launches from spaceX API  and cancelling internal bookings that are scheduled on the same launchpad for the same day.
When a cancellation occurs, an event is dispatched to an event broker to notify other services.

Additionally, users can cancel their bookings by sending their ticket ID to the API.

Note:
- Due to time constraints, not all infrastructure package tests have been fully covered. However, some tests have been implemented to showcase specific scenarios.
- Although the challenge mentions 2049 as the present year, i didn't restrict the API to only accept bookings after this date because the spaceX API is returning upcoming launches only in 2022;

## Installation

### Prerequisites
- Go 1.23+
- Docker (optional, for containerization)
- Buf (optional, for building protobufs)
- golang-migrate (optional, for interact with database migrations)
- mockery (optional, for generating the mocks)

### Steps
1. Install optional dependencies:
    - Docker: Tested using [colima](https://github.com/abiosoft/colima);
    - [golang-migrate](https://github.com/golang-migrate/migrate);
    - [buf](https://buf.build/docs/installation)
    - [mockery](https://github.com/vektra/mockery)

2. ``` go mod download ```

## Usage

#### Running Locally:
To run the service locally, one needs to run the following command:
``` go run cmd/spacet/main.go -config=./config/config.yaml ```
The `-config` flag defines the path of the configuration file. The default config file may produce errors due to missing values. To resolve this, you can either:
1. Update the configuration file with the required values.
2. use [environment variables](#Environment_Variables)

##### minimal environment with the default config:
```bash 
export PG_DSN="host=localhost port=5432 user=postgres dbname=spacet-db password=spacetpw sslmode=disable"
export PUBSUB_EMULATOR_HOST="localhost:8681"
            
```

#### Running via Docker Compose
Running ``` docker compose up --detach ``` will bring up the following services:
- PostgreSQL: The database service.
- Spacet Microservice
- PubSub Emulator: Acts as an event broker.
- Echo Service: Consumes and echoes PubSub messages.
- Publish Service: Ensures the PubSub emulator is functioning correctly by sending messages.


#### Makefile
The makefile provides the following commands:
- `make proto`- Updates the project's protobufs.
- ` make mocks `- Generates mocks using Mockery.
- ` make docker-compose `- Runs Docker Compose in detached mode.
- ` make test ` - Runs the tests.


## Environment_variables

| Variable Name | Description                            
|---------------|----------------------------------------
| `APP_NAME`        | The name of the application.    
| `APP_VERSION`        | The version of the application.    
| `LOG_LEVEL`   | The log level for the application
| `HTTP_PORT`  | The port on which the HTTP server runs.   
| `GRPC_PORT`   | The port on which the GRPC server runs. 
| `PG_POOL_MAX`      | The maximum number of connections in the Postgres pool. 
| `PG_DSN`      | The Data Source Name (DSN) for connecting to Postgres. 
| `GIN_MODE`      | Sets the gin mode (http server). ex: "release"
| `ORCHESTRATOR_INTERVAL` | Manages the schedule interval that syncs spacex launches
| `PUBSUB_EMULATOR_HOST`      | Sets the pubsub emulator host. 
| `NOTIFICATIONS_BATCH_SIZE_MAX`      | The maximum batch size for processing notifications. 
| `NOTIFICATIONS_INTERVAL`      | The interval (in seconds) for sending notifications. 
| `PUBSUB_ENABLED`      | Flag to enable or disable Pub/Sub functionality. 
| `PUBSUB_PROJECT_ID`      | The Google Cloud project ID for Pub/Sub. 
| `PUBSUB_LAUNCHES_TOPIC`      | The Pub/Sub topic for launches-related events. 



## API

This service provides both an HTTP server and a GRPC server. The protocol buffer definitions are located in `/protos/spacet.proto`.

### HTTP server
The HTTP server offers monitoring endpoints and acts as a gateway to the GRPC server using the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

Default port: `8080`

The monitoring routes are:

`GET /healthz` - Returns a status code 200 when the api is responsive

`GET /readiness` - Returns a status code 200 and a json body when the api is ready

`GET /liveness` - Monitors dependencies and returns 200 or 500 status codes.

The gateway routes are served in the subpath /v1 

Check ```/protos/spacet.proto``` or ```/gen/proto/openapiv2/spacet.swagger.json``` 

### GRPC server
The GRPC server uses protovalidate to validate proto message fields. Note that server reflection is not enabled, so you may need to import the validate proto manually. It can be downloaded from (https://github.com/bufbuild/protovalidate.git)

Default port: `8081`

### Postman Collection
A Postman collection is included to simplify testing. However, it does not include any automation.

directory: `/examples/spacet.postman_collection.json`

## Implementation

### Naming
There are 4 main Entities in the domain model. 
LaunchPads, Launches, Bookings and Tickets. A Ticket is essentially a booked launch. 

### SpaceX API
Inside /pkg/spacexapi there is an implementation of the spacexclient v4 api (only features launchpads, launches and landpads)

### LaunchPad Fetching
The launchpads are fetched every 30 days (currently upon application restart, but the intent was to create a schedule)

### Launches Sync
The sync between spaceX launches and internal launches is done using a scheduled. The frequency interval can be configured.
This routine will cancel internal launches that collide with the spaceX schedule. In the future we can send a notification when bookings are cancelled.

### Distributed Locks
The API uses distributed locks (using pg_try_advisory_lock), please check /internal/app/sync/commands.go, namely on SyncIfNecessary implementation, for more details.

### PubSub connection
to use the pubsub emulator the environment variable `PUBSUB_EMULATOR_HOST` needs to be set.For production, refer to the official Google documentation on setting up credentials and configuring the required environment variables.

### Outbox Pattern
We use this pattern to ensure that notifications are sent asynchronously while keeping the dependencies minimal. 
This aligns well with the SOLID and CQRS principles used here.
The outbox table keeps the events after they are successfully sent, this can be changed to remove right after the publish is made or by creating another async process to cleanup the table after some time.
Note that this implementation may result in message duplication in the message broker.

### Cursor Based Pagination
The list endpoint implements cursor-based pagination.

### Design Patterns
This project tries to follow SOLID, CQRS and clean architecture as much as possible. While there are opportunities for optimization, time constraints prevented full exploration.


### Folder structure
The `/pkg` folder provides Interfaces to interact with postgres, pubsub, logger and the servers.

The `/migrations` folder contains the sql migrations

The `/gen` folder contains the generated mocks, go proto packages, swagger definitions...

The `/examples` folder intent was to have usage examples of the api. It also contains some dockerfiles that act as an example for publishing and consuming pubsub events

The `/config` folder contains the base config and config structure of the application

The `/cmd/spacet` is the entry point of the microservice

Finally, the `/internal` folder contains the app implementation:
 - `/app`  contains the service/usecase layer implementation
 - `/controller` contains the controller/handlers layer (grpc and http apis)
 - `/domain` contains the application domain definitions (repo interfaces, entities, models)
 - `/infra` contains the infrastructure layer with repositories, notifications, outbox implementation ....



### DB Migrations
 ` migrate create -ext sql -dir migrations/postgresql -seq filename `
