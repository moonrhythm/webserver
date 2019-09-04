FROM gcr.io/moonrhythm-containers/golang:1.13.0-alpine as build

ENV CGO_ENABLED=0
WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o webserver -ldflags '-w -s' .

FROM scratch

COPY --from=build /workspace/webserver /webserver

WORKDIR /www

ENTRYPOINT ["/webserver"]
