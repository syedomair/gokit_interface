package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	var (
		httpAddr = flag.String("http.addr", ":"+port, "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s Service
	{
		s = DBService()
	}

	var h http.Handler
	{
		h = MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	logger.Log("transport", "HTTP", "addr", *httpAddr)
	//http.ListenAndServe(*httpAddr, appHandler(handler))
	http.ListenAndServe(*httpAddr, h)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

type appHandler func(http.ResponseWriter, *http.Request)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(w, r)
}
