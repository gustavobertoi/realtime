## Build stage

FROM golang:1.21 as build

WORKDIR /app

COPY . /app

RUN go mod download
RUN CGO_ENABLED=0 go build -o bin/realtime .

## Runnable container

FROM gcr.io/distroless/base-debian11 as runnable

COPY --from=build /app/bin/realtime /

EXPOSE 8080

USER nonroot:nonroot

CMD ["/realtime"]
