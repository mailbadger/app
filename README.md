[![Build Status](https://travis-ci.org/news-maily/app.svg?branch=master)](https://travis-ci.org/news-maily/app)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/FilipNikolovski/news-maily/master/LICENSE.md)

# news-maily

Self hosted newsletter mail system written in go.

# Development setup


This application consists of an Rest API written in Go and a dashboard application which is written in React. The whole app resides in this single repository.

The application depends on several tools and services:
    
    - go
    - MySQL
    - NSQ
    - yarn
    - go-bindata
    - Docker and docker-compose (optional)

1. `go-bindata` is used to generate the DB migration assets (sql files)

```
go get -u github.com/go-bindata/go-bindata/...
```

2. Run `make gen` to generate the migration assets.

3. Run `make build` to build the executable files, the files will be located in the `release` folder.

4. Run MySQL and NSQ services (see the docker-compose.yaml file).

5. (Optional) if you run the application using docker-compose, you'll need to copy the contents of .example.env to .env so the environment variables will be read by the application.

## Starting the application

