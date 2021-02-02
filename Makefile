.PHONY: build

UID=$$(id -u)
GID=$$(id -g)

all: build run

build:
	go build -o ./build/godmin -v ./cmd/main.go

run:
	./build/godmin

test: migrate_test_up
	GODMIN_ENV=test && go test -v -race -timeout 30s ./...

migrate_create:
	docker run -u ${UID}:${GID} -v ${PWD}/migrations:/migrations migrate/migrate \
		create -ext sql -dir /migrations $(name)

migrate= \
	docker run -u ${UID}:${GID} -v ${PWD}/migrations:/migrations --link $(1) --net godmin_default migrate/migrate \
		-path=/migrations/ \
		-database postgres://godmin:password@$(1):5432/$(1)?sslmode=disable \
		$(2)

migrate_up:
		$(call migrate,godamin_db_dev,up)

migrate_down:
		$(call migrate,godamin_db_dev,down)

migrate_test_up:
		$(call migrate,godamin_db_test,up)

migrate_test_down:
		$(call migrate,godamin_db_test,down)
