package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "[HttpServer] ", log.LstdFlags | log.Lshortfile)


func IndexHandler(w http.ResponseWriter, r *http.Request) {

	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}

	os.Setenv("VERSION", "0.0.4")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)

	ip := getRemoteClientIp(r)
	code := http.StatusAccepted
	logger.Printf("ip: %s, status: %d \n", ip, code)

	w.Write([]byte("<h1>hello, world</h1>"))
}

func getRemoteClientIp(r * http.Request) string {
	xRealIp := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xRealIp == "" && xForwardedFor == "" {
		idx := strings.LastIndex(r.RemoteAddr, ":")
		return r.RemoteAddr[:idx]
	}

	if xForwardedFor != "" {
		xForwardedFor := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(xForwardedFor[0])
	}

	return xRealIp
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>healthz test</h1>"))
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/healthz", healthzHandler)

	logger.Print("server listening on 127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		logger.Fatal("failed to start server")
	}
}
