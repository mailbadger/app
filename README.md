[![Go Report Card](https://goreportcard.com/badge/github.com/mailbadger/app)](https://goreportcard.com/report/github.com/mailbadger/app)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mailbadger/app)
[![codecov](https://codecov.io/gh/mailbadger/app/branch/master/graph/badge.svg)](https://codecov.io/gh/mailbadger/app)

![Mailbadger Logo](https://github.com/mailbadger/app/blob/assets/Mailbadger_MascotWordMarkOutline_Black.png?raw=true "Mailbadger Logo")

Self hosted newsletter mail system written in go.

# Installation instructions

This application consists of an Rest API written in Go and a dashboard application which is written in React. The whole app resides in this single repository.

The application depends on several tools and services:

    - go
    - MySQL (or sqlite)
    - NSQ
    - yarn
    - statik
    - Docker and docker-compose (optional)

1. `statik` is used to generate the DB migration assets (sql files)

```
go get github.com/rakyll/statik
```

2. Run `make driver=mysql gen` to generate the migration assets (run with driver=sqlite3 for testing).

3. Run `make build` to build the executable files, the files will be located in the `bin` folder.

4. Run MySQL and NSQ services (see the docker-compose.yaml file).

5. See the `.example.env` file to see which env variables should be set for the application to run.

## Starting the application
