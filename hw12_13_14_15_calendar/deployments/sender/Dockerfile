FROM golang:1.21 as build

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS=linux

COPY . .

RUN go build -o bin/calendar_sender ./cmd/calendar-sender/

FROM busybox

COPY --from=build /app/bin/calendar_sender .
COPY --from=build /app/configs/config.yml .

CMD ["/calendar_sender", "-config", "/config.yml"]
