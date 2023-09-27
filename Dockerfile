FROM golang:1.21 as build

WORKDIR /app

COPY . /app/

RUN go mod download
RUN go build -o /bin/realtime cmd/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/bin/realtime .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./realtime"]
