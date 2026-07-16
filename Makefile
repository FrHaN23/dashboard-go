# Variables
BINARY_NAME=main
go_bin ?= go
DB_SCRIPT_PATH=./script/db/init-db.go
SERVER_MAIN_PATH=./main.go

# Default target
.DEFAULT_GOAL := help

dep:
	@go mod tidy
	@go mod vendor

## run: Menjalankan aplikasi/server Go secara langsung
run:
	CGO_ENABLED=1 $(go_bin) run $(SERVER_MAIN_PATH)

## build: Melakukan kompilasi program menjadi file binary di dalam folder bin/
build:
	CGO_ENABLED=1 $(go_bin) build -o bin/${BINARY_NAME} $(SERVER_MAIN_PATH)

## db-init: Menjalankan inisialisasi tabel basis data dan melakukan seeding data awal
db-init:
	@echo "Running database initialization and seeding..."
	go run $(DB_SCRIPT_PATH)

## help: Menampilkan daftar perintah yang tersedia di dalam Makefile ini
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
