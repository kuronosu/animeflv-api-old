FROM golang:alpine AS build

COPY ["go.mod", "go.sum", "/go/src/app/"]

WORKDIR /go/src/app

RUN go mod download

COPY [".", "."]

RUN go build -o /go/bin/deguvon-api ./cmd/deguvon-api \
    && go build -o /go/bin/deguvon-cli ./cmd/deguvon-cli

FROM alpine:latest

WORKDIR /go/bin

COPY --from=build ["/go/bin/deguvon-api", "/go/bin/deguvon-cli", "/go/bin/"]
COPY --from=build ["/go/src/app/res", "/go/bin/res/"]

EXPOSE 8080

CMD ["/go/bin/deguvon-api"]
