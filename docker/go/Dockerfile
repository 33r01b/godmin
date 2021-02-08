# Stage 1. Install
FROM golang:1.15.8 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Stage 2. Build
FROM golang:1.15.8 as builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 1000 godmin

RUN mkdir -p /godmin
ADD . /godmin
WORKDIR /godmin

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/godmin -v ./cmd/godmin/main.go

# Stage 3. Run
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER godmin

COPY --from=builder /godmin/bin/godmin /godmin
COPY --from=builder /godmin/.env /.env

CMD ["/godmin"]