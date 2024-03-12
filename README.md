# Konsep Domain Driven Design

from - https://leravio.com/blog/domain-driven-design-dalam-golang/

Domain-Driven Design Golang

- bounded contexts: setiap aplikasi memiliki masing-masing (domain, repository, service, dan controller yang berbeda-beda). membantu mendefinisikan batasan dan interaksi setiap domain
- aggregates: kumpulan objek domain yang saling terkait, menjamin enkapsulasi ke keadaan internal mereka. (konsep struct, interface, method)
- Entities: menggunakan konsep struct untuk identifikasi data model representasi tabel di database, atau data transfer objek
- Ubiquitous Language: penggunaan bahasa yang konsisten di antara semua pemangku kepentingan dalam proyek (penamaan variable, fungsi, method, dll)

## Architecture

1. Domain -> data model dari representasi tabel, dan data model request & response
2. Repository -> data akses layer
3. Service -> adalah bisnis logic, seperti transaction db, add request to data model tabel, panggil repository, atau penegcekan data
4. Event -> untuk menangani request eksternal sistem, seperti (HTTP, gRPC, Messaging, etc)
5. Application -> mengkoordinasi interaksi antara lapisan presentasi (API) ke service, menerima input, dan mengoordinasikan operasi (service) yang diperlukan
6. Infrastructure -> menagani aspek teknis seperti database, integrasi eksternal lainnya


## Tech Stack

- Golang : https://github.com/golang/go
- MySQL (Database) : https://github.com/mysql/mysql-server

## Framework & Library

- GoFiber (HTTP Framework) : https://github.com/gofiber/fiber
- GORM (ORM) : https://github.com/go-gorm/gorm
- Viper (Configuration) : https://github.com/spf13/viper
- Golang Migrate (Database Migration) : https://github.com/golang-migrate/migrate
- Go Playground Validator (Validation) : https://github.com/go-playground/validator
- Logrus (Logger) : https://github.com/sirupsen/logrus

## Configuration

Semua data konfigurasi ada di file `config.json`.

## API Spec

Semua API Spec ada di folder api.

## Database Migration

Semua database migration ada di folder db/migrations.

### Install Migrate
```shell
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Create Migration

```shell
migrate create -ext sql -dir db/migrations create_table_xxx
```

### Run Migration

```shell
migrate -database "mysql://username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up
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