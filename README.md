# simple-calendar

## Local run:
Pull this project: `git pull https://github.com/farir1408/simple-calendar`

Before running service, change `LOCAL_GOOSE` variable in Makefile on your local OS.
Then run `make build` and run service with `--run-local` flag.

After running service start grpc listeners on `CALENDAR_API_PORT` and http listeners on `CALENDAR_API_DEBUG_PORT`.
Openapi documentation is available on `localhost:CALENDAR_API_DEBUG_PORT/swaggerui/` end-point. 

## Customize proto definition
For customize api:
 - change api/calendar.proto file.
 - run `make install`
 - run `make generate`
 
For more information about grpc-gateway see: https://github.com/grpc-ecosystem/grpc-gateway