FROM golang:1.23.4 as build

ENV CGO_ENABLED=0
WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o webserver -ldflags '-w -s' .

FROM gcr.io/distroless/static

COPY --from=build --link /workspace/webserver /app/webserver

WORKDIR /www

ENTRYPOINT ["/app/webserver"]
