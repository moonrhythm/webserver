package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/moonrhythm/webstatic/v4"
)

var (
	indexFile          = flag.String("index", "index.html", "index file")
	indexCacheControl  = flag.String("index.cache-control", "no-cache", "index cache control")
	notFoundFile       = flag.String("notfound", "404.html", "404 file")
	spa                = flag.Bool("spa", true, "spa mode")
	port               = flag.Int("port", 8080, "http port (override by env PORT)")
	assetsCacheControl = flag.String("asset.cache-control", "public, max-age=3600", "assets cache control")
	dir                = flag.String("dir", ".", "serve dir")
)

var (
	serveDir   http.FileSystem
	fileServer http.Handler
)

func main() {
	flag.Parse()

	envPort := os.Getenv("PORT")
	if envPort != "" {
		p, _ := strconv.Atoi(envPort)
		if p > 0 {
			*port = p
		}
	}

	serveDir = http.Dir(*dir)
	fileServer = http.FileServer(serveDir)

	log.Printf("start web server on %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *port), &webstatic.Handler{
		FileSystem:   serveDir,
		CacheControl: *assetsCacheControl,
		Fallback:     http.HandlerFunc(index),
	}))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if filepath.Ext(r.URL.Path) == "" && tryServeHTML(w, r) {
			return
		}

		if !*spa {
			if *notFoundFile != "" {
				http.ServeFile(w, r, *notFoundFile)
				return
			}
			http.NotFound(w, r)
			return
		}
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
