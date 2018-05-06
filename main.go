package main

import (
	"fmt"
	"html"
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
	httpAddr := os.Getenv("NOMAD_PORT_http")
	if httpAddr == "" {
		log.Fatal("NOMAD_PORT_http must be set and non-empty")
	}

	pathCert := os.Getenv("PATH_CERT")
	pathKey := os.Getenv("PATH_KEY")

	listenAddr := fmt.Sprintf(":%v", httpAddr)
	log.Printf("Listening on [%v]...\n", listenAddr)

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/healthz", handleHealthzRequest)

	var err error
	if pathCert != "" && pathKey != "" {
		err = http.ListenAndServeTLS(listenAddr, pathCert, pathKey, nil)
	} else {
		err = http.ListenAndServe(listenAddr, nil)
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return 0
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request [%v]...\n", html.EscapeString(r.URL.Path))
	w.Write([]byte(ResponseOk))
	return
}

func handleHealthzRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(ResponseOk))
	return
}
