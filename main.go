package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr               = flag.String("addr", "0.0.0.0:3000", "TCP address to listen to")
	dir                = flag.String("dir", "/var/www/html", "Directory to serve static files from")
	generateIndexPages = flag.Bool("generateIndexPages", false, "Whether to generate directory index pages")
	spaMode            = flag.Bool("spaMode", false, "Single Page Application Mode")
)

func main() {
	const API_CONTENT_TYPE = "application/vnd.brickdoc.app-engine+json; charset=utf8"
	flag.Parse()

	// Setup Fasthttp file server handler

	notFoundHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType(API_CONTENT_TYPE)
		ctx.SetBodyString("{\"error\": \"Resource Not found\", \"handler\":\"WrapDrive Static Server (Brickdoc App Engine)\"}")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: *generateIndexPages,
		PathNotFound:       notFoundHandler,
		// Gzip and brotil compress by ingress or load balancing.
		Compress:        false,
		AcceptByteRange: true,
	}

	if *spaMode {
		fs.PathRewrite = spaPathRewrite()
	}

	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {

		if string(ctx.Path()) == "/healthz" {
			ctx.SetContentType(API_CONTENT_TYPE)
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
