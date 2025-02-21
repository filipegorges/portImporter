# Port Importer

Imports ports from a JSON file passed in as argument to the application.

### Expected format:
```json
{
  "AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "Abu Dhabi",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  },
}
```

## Running the application

For ease of use, this application can be run through docker compose:

```shell
docker compose up -d
```

## Checking the created data

For ease of use, `mongo-express` has been added to the docker compose, with the application running:

http://localhost:8081/db/portImporter/
* user: admin
* pass: pass

NOTE: this is for local verification only

The data will be created within the `portImporter` collection.

## Running the tests

Within the project's root, run:

```shell
go test ./...
```

NOTE: In the future, we could add a Makefile to centralize all these commands

## Running the linter

With docker up and running, run:

```shell
docker-compose run --rm linter
```

NOTE: no output means no errors found!
