package main

import (
	"encoding/json"
	"flag"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	var (
		httpAddr = flag.String("http.addr", ":"+port, "HTTP listen address for CloudPit")
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

	var sh http.Handler
	{
		sh = securityMiddleware(s, h)
	}

	logger.Log("transport", "HTTP", "addr", *httpAddr)
	http.ListenAndServe(*httpAddr, sh)
}

func securityMiddleware(s Service, handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-key")
		jwtToken := r.Header.Get("x-jwt")
		urlPath := r.URL.Path

		rtn := s.AuthProvider(apiKey, jwtToken, urlPath)
		if rtn == nil {
			handle.ServeHTTP(w, r)

		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(rtn)
		}
	})
}
