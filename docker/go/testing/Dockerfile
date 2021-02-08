FROM golangci/golangci-lint:v1.36-alpine

RUN mkdir /godmin
ADD . /godmin
WORKDIR /godmin

ENV GODMIN_ENV=test

RUN golangci-lint run --issues-exit-code=1 --deadline=600s ./...

RUN go test -v ./...