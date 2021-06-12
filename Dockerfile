FROM golang:1.16.5 as build

ENV CGO_ENABLED=0
WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o webserver -ldflags '-w -s' .

FROM gcr.io/distroless/static

COPY --from=build /workspace/webserver /app/webserver

WORKDIR /www

ENTRYPOINT ["/app/webserver"]
