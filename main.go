package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr               = flag.String("addr", "0.0.0.0:8080", "TCP address to listen to")
	dir                = flag.String("dir", "/var/www/html", "Directory to serve static files from")
	generateIndexPages = flag.Bool("generateIndexPages", false, "Whether to generate directory index pages")
	spaMode            = flag.Bool("spaMode", false, "Single Page Application Mode")
)

func main() {
	flag.Parse()

	// Setup FS handler

	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: *generateIndexPages,
		// gzip compress by nginx ingress or load balancing.
		Compress:        false,
		AcceptByteRange: true,
	}

	if *spaMode {
		fs.PathRewrite = spaPathRewrite()
	}

	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {

		if string(ctx.Path()) == "/healthz" {
			ctx.SetContentType("text/plain; charset=utf8")
			ctx.SetBodyString("{\"alive\": true}")
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
		}
		fsHandler(ctx)
	}

	// Start server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		go func() {
			if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
				log.Fatalf("error in ListenAndServe: %s", err)
			}
		}()
	}

	log.Printf("Serving files from directory %q", *dir)
	log.Printf("See stats at http://%s/healthz", *addr)

	// Wait forever./
	select {}
}

func spaPathRewrite() fasthttp.PathRewriteFunc {
	return func(ctx *fasthttp.RequestCtx) []byte {
		log.Print(ctx.Path())
		if contains(ctx.Path(), byte('.')) {
			return ctx.Path()
		} else {
			return []byte("/index.html")
		}
	}
}

func contains(arr []byte, byt byte) bool {
	for _, a := range arr {
		if a == byt {
			return true
		}
	}
	return false
}