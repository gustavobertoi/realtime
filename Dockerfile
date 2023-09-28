FROM golang:1.21 as build

WORKDIR /app

COPY . /app

RUN go mod download
RUN go build -o /app/bin/realtime /app/cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=build /app/bin/realtime .

EXPOSE 8080

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser:appgroup

ENTRYPOINT ["./realtime"]
