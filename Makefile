.PHONY: build

all: build run

build:
	go build -o ./build/godmin -v ./cmd/main.go

run:
	sudo ./build/godmin