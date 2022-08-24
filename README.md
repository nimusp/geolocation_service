# Geolocation service

FindHotel Coding Challenge

## How to run

First run the data import to prepare the database, this should take about a minute

`$ docker-compose up importer`

run the service

`$ docker-compose up gateway`

## Description

Repository contains the importer `/cmd/importer` and the gateway `/cmd/gateway` components.

The importer parses the passed .csv file, sanitizes it and inserts it into the database using batch insertion. If the database contains dublicate IP addreses, the importer rewrites it with the new data.

The gateway has a single endpoint `/location/{ipAddress}` to return geo data on the passed IP.

## Setting up the environment

- IMPORT_FROM — path to .csv file for import
- IMPORT_BATCH_SIZE — the batch size to be inserted into the database (default 1000)
- DB_HOST — database hostname
- DB_PORT — database port
- DB_USER — database user name
- DB_PASSWORD — database user password
- DB_NAME — database name
- GATEWAY_PORT — service port (default :8888)
- GATEWAY_TIMEOUT_SECONDS — service timeout (default 5 second)

## Makefile targets

- help — show this help
- deps — install required dependencies
- fmt — format source files of the project
- test — run all the tests
- mocks — rebuild all the mocks
- import — run the importer with docker-compose
- gateway — run the service with docker-compose
