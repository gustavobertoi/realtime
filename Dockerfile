## Build stage

FROM golang:1.22 as build

WORKDIR /app

COPY . /app

RUN go mod download
RUN CGO_ENABLED=0 go build -o bin/realtime /app/cmd/api/main.go

## Runnable container

FROM gcr.io/distroless/base-debian11 as runnable

COPY --from=build /app/bin/realtime /
COPY --from=build /app/web /web

EXPOSE 8080

CMD ["/realtime"]
