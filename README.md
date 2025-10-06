# Golang Clean Architecture Template

## Description

This is golang clean architecture template.

## Tech Stack

- Golang : https://github.com/golang/go
- PostgreSQL (Database) : https://github.com/lib/pq

## Framework & Library

- Go-Chi (HTTP Router) : https://github.com/go-chi/chi
- GORM (ORM) : https://github.com/go-gorm/gorm
- Godotenv (Configuration) : https://github.com/joho/godotenv
- Golang Migrate (Database Migration) : https://github.com/golang-migrate/migrate
- Go Playground Validator (Validation) : https://github.com/go-playground/validator
- Logrus (Logger) : https://github.com/sirupsen/logrus

## Configuration

Example configuration is in `.env.example` file.
Add two configuration files, namely “.env” and “.env.test”. 

## API Spec

All API Spec is in `api` folder.

## Database Migration

All database migration is in `db/migrations` folder.

### Create Migration

```shell
migrate create -ext sql -dir db/migrations -seq create_xxx_table
```

### Run Migration

```shell
migrate -database postgres://postgres:password@localhost:5432/example?sslmode=disable -path db/migrations up
```

## Run Application

### Run unit test

```bash
go test -v ./test/
```

### Run web server

```bash
go run cmd/web/main.go
```
