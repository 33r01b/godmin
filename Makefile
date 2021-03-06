.PHONY: build

UID=$$(id -u)
GID=$$(id -g)

all: build run

build:
	go build -o ./bin/godmin -v ./cmd/godmin/main.go

image:
	docker build -t 33r01b/godmin -f docker/go/Dockerfile .

run:
	./bin/godmin

test: migrate_test_up
	docker build --network host -t 33r01b/godmin-test -f docker/go/testing/Dockerfile .

migrate_create:
	docker run -u ${UID}:${GID} -v ${PWD}/migrations:/migrations migrate/migrate \
		create -ext sql -dir /migrations $(name)

migrate= \
	docker run -u ${UID}:${GID} -v ${PWD}/migrations:/migrations --link $(1) --net godmin_default migrate/migrate \
		-path=/migrations/ \
		-database postgres://godmin:password@$(1):5432/$(1)?sslmode=disable \
		$(2)

migrate_up:
		$(call migrate,godmin_db_dev,up)

migrate_down:
		$(call migrate,godmin_db_dev,down)

migrate_test_up:
		$(call migrate,godmin_db_test,up)

migrate_test_down:
		$(call migrate,godmin_db_test,down)
