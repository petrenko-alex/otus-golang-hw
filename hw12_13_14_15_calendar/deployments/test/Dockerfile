FROM golang:1.21 as build

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS=linux

COPY . .

RUN go test -o bin/server.test -c ./internal/server/server_test.go

FROM busybox

ENV CONFIG="./config.test.yml"

COPY --from=build /app/bin/server.test .
COPY --from=build /app/configs/config.test.yml .

CMD ["/server.test"]
