FROM golang:1.21 as build

WORKDIR /app

COPY . /

RUN go mod download
RUN go build -o /bin/realtime cmd/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /bin/realtime .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./realtime"]
