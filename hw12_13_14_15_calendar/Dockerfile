FROM golang:1.21 as build

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS=linux

COPY . .

RUN go build -o bin/calendar ./cmd/calendar/
RUN go build -o bin/migrate ./cmd/calendar-migrations/

FROM busybox

EXPOSE 4042
EXPOSE 4242

COPY --from=build /app/bin/calendar .
COPY --from=build /app/bin/migrate .
COPY --from=build /app/migrations ./migrations/
COPY --from=build /app/configs/config.yml .

CMD ["sh", "-c", "/migrate -config=./config.yml up && /calendar -config=./config.yml"]
