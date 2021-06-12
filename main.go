package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/moonrhythm/webstatic"
)

var (
	indexFile          = flag.String("index", "index.html", "index file")
	indexCacheControl  = flag.String("index.cache-control", "no-cache", "index cache control")
	notFoundFile       = flag.String("notfound", "404.html", "404 file")
	spa                = flag.Bool("spa", true, "spa mode")
	port               = flag.Int("port", 8080, "http port (override by env PORT)")
	assetsCacheControl = flag.String("asset.cache-control", "public, max-age=3600", "assets cache control")
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

	log.Printf("start web server on %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), webstatic.New(webstatic.Config{
		CacheControl: *assetsCacheControl,
		Fallback:     http.HandlerFunc(index),
	})))
}

func index(w http.ResponseWriter, r *http.Request) {
	if !*spa && r.URL.Path != "/" {
		if *notFoundFile != "" {
			http.ServeFile(w, r, *notFoundFile)
			return
		}
		http.NotFound(w, r)
		return
	}

	if *indexCacheControl != "" {
		w.Header().Set("Cache-Control", *indexCacheControl)
	}
	http.ServeFile(w, r, *indexFile)
}
