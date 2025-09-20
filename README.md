### Create Migration

```shell
migrate create -ext sql -dir db/migrations -seq create_xxx_table

### Run Migration

```shell
migrate -database postgres://postgres:password@localhost:5432/example?sslmode=disable -path db/migrations up