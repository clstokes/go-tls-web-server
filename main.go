package main

import (
	"flag"
	"github.com/armon/go-metrics"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	ResponseOk = "{\"status\":\"ok\"}"
)

type ServerOpts struct {
	listenAddr       string
	clientCert       string
	clientKey        string
	maxCrashDuration int
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	server := &ServerOpts{}
	parseArgs(server)

	if server.maxCrashDuration != 0 {
		setupCrashRoutine(server.maxCrashDuration)
	}

	log.Printf("Listening on [%v]...\n", server.listenAddr)

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/healthz", handleHealthzRequest)

	var err error
	if server.clientCert != "" && server.clientKey != "" {
		err = http.ListenAndServeTLS(server.listenAddr, server.clientCert, server.clientKey, nil)
	} else {
		err = http.ListenAndServe(server.listenAddr, nil)
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

func setupCrashRoutine(maxCrashDuration int) {
	rand.Seed(time.Now().Unix())
	crashDuration := rand.Intn(maxCrashDuration)

	log.Printf("Crashing in [%v] seconds...", crashDuration)
	go func() {
		time.Sleep(time.Duration(crashDuration) * time.Second)
		log.Fatal("Crashing...")
	}()
}

func parseArgs(server *ServerOpts) {
	flag.StringVar(&server.listenAddr, "listen-address", ":8080", "address to listen on")
	flag.StringVar(&server.clientCert, "client-cert", "", "path to the CA certificate")
	flag.StringVar(&server.clientKey, "client-key", "", "path to the CA private key")
	flag.IntVar(&server.maxCrashDuration, "crash", 0, "maximum duration to wait before crashing")
	flag.Parse()
}
