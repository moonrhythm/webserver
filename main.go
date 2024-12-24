package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/moonrhythm/parapet"
	"github.com/moonrhythm/parapet/pkg/logger"
	"github.com/moonrhythm/webstatic/v4"
)

var (
	indexFile          = flag.String("index", "index.html", "index file")
	indexCacheControl  = flag.String("index.cache-control", "no-cache", "index cache control")
	notFoundFile       = flag.String("notfound", "404.html", "404 file")
	spa                = flag.Bool("spa", true, "spa mode")
	port               = flag.String("port", "8080", "http port (override by env PORT)")
	bindIP             = flag.String("bind-ip", "", "ip to listen")
	assetsCacheControl = flag.String("asset.cache-control", "public, max-age=3600", "assets cache control")
	dir                = flag.String("dir", ".", "serve dir")
)

var (
	serveDir       http.FileSystem
	fileServer     http.Handler
	notFoundBuffer []byte
)

func main() {
	flag.Parse()

	envPort := os.Getenv("PORT")
	if envPort != "" {
		*port = envPort
	}

	serveDir = http.Dir(*dir)
	fileServer = http.FileServer(serveDir)
	notFoundBuffer = []byte("404 page not found")
	loadBuffer(*notFoundFile, &notFoundBuffer)

	srv := parapet.NewBackend()
	srv.Addr = net.JoinHostPort(*bindIP, *port)
	srv.H2C = true
	srv.Handler = &webstatic.Handler{
		FileSystem:   serveDir,
		CacheControl: *assetsCacheControl,
		Fallback:     http.HandlerFunc(index),
	}

	srv.Use(logger.Stdout())
	srv.Use(parapet.MiddlewareFunc(securityHeaders))

	log.Printf("start web server on %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if filepath.Ext(r.URL.Path) == "" && tryServeHTML(w, r) {
			return
		}

		if !*spa {
			serveNotFound(w, r)
			return
		}

		serveIndexFallback(w, r)
		return
	}

	setCacheControl(w)
	http.ServeFile(w, r, *indexFile)
}

func setCacheControl(w http.ResponseWriter) {
	if *indexCacheControl == "" {
		return
	}
	w.Header().Set("Cache-Control", *indexCacheControl)
}

func tryServeHTML(w http.ResponseWriter, r *http.Request) (served bool) {
	htmlPath := r.URL.Path + ".html"

	fs, err := serveDir.Open(htmlPath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return true
	}
	fs.Close()

	r.URL.Path = htmlPath
	setCacheControl(w)
	fileServer.ServeHTTP(w, r)
	return true
}

func loadBuffer(fn string, buf *[]byte) {
	if fn == "" {
		return
	}

	tmpBuff, err := os.ReadFile(fn)
	if err != nil {
		return
	}
	*buf = tmpBuff
}

func serveIndexFallback(w http.ResponseWriter, r *http.Request) {
	// prevent new files that were not fully serve yet to cache on cdn
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.ServeFile(w, r, *indexFile)
}

func serveNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	io.Copy(w, bytes.NewReader(notFoundBuffer))
}

func securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		h.ServeHTTP(w, r)
	})
}
