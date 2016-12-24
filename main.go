package main

import (
	"log"
	"net/http"
	"os"
)

const (
	ResponseOk = "{\"status\":\"ok\"}\n"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	httpAddr := os.Getenv("NOMAD_ADDR_http")
	if httpAddr == "" {
		log.Fatal("NOMAD_ADDR_http must be set and non-empty")
	}

	pathCert := os.Getenv("PATH_CERT")
	if pathCert == "" {
		log.Fatal("PATH_CERT must be set and non-empty")
	}

	pathKey := os.Getenv("PATH_KEY")
	if pathKey == "" {
		log.Fatal("PATH_KEY must be set and non-empty")
	}

	log.Printf("Listening on [%v]...\n", httpAddr)

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/healthz", handleHealthzRequest)
	err := http.ListenAndServeTLS(httpAddr, pathCert, pathKey, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return 0
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(ResponseOk))
	return
}

func handleHealthzRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(ResponseOk))
	return
}
