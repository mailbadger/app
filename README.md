[![Go Report Card](https://goreportcard.com/badge/github.com/news-maily/app)](https://goreportcard.com/report/github.com/news-maily/app)
[![Go 1.12](https://img.shields.io/badge/go-1.12-9cf.svg)](https://golang.org/dl/)

# news-maily

Self hosted newsletter mail system written in go.

# Installation instructions

This application consists of an Rest API written in Go and a dashboard application which is written in React. The whole app resides in this single repository.

The application depends on several tools and services:
    
    - go
    - MySQL (or sqlite)
    - NSQ
    - yarn
    - go-bindata
    - Docker and docker-compose (optional)

1. `go-bindata` is used to generate the DB migration assets (sql files)

```
go get -u github.com/go-bindata/go-bindata/...
```

2. Run `make gen` to generate the migration assets.

3. Run `make build` to build the executable files, the files will be located in the `bin` folder.

4. Run MySQL and NSQ services (see the docker-compose.yaml file).

5. See the `.example.env` file to see which env variables should be set for the application to run.

## Starting the application