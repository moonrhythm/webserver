# webserver

Small Web Server to serve SPA and static files

## Usage

### Dockerfile

```Dockerfile
FROM gcr.io/moonrhythm-containers/webserver

ADD . .
```

### Build

```sh
docker build -t web .
```

### Run

```sh
docker run --rm -p 8080:8080 web
```
