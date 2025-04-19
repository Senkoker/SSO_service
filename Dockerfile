FROM golang:1.23.4-alpine AS Builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ./go.mod ./

COPY ./go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app/migration ./main/migrations/main.go

RUN go build -o ./bin/app/cmd ./main/main.go

FROM alpine

COPY --from=builder /usr/local/src/bin/app/migration ./

COPY --from=builder /usr/local/src/bin/app/cmd ./

COPY ./skript.sh ./

COPY ./internal/migration_db ./migration_db

CMD ["./skript.sh"]