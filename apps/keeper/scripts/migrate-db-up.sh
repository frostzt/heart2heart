#!/bin/sh
migrate -database postgres://postgres:postgres@localhost:5432/keeper?sslmode=disable -path db/migrations up
