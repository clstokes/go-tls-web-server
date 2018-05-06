package main

import (
	"flag"
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

func main() {
	os.Exit(realMain())
}

func realMain() int {
	listenAddr := flag.String("listen-address", ":8080", "address to listen on")
	clientCert := flag.String("client-cert", "", "path to the CA certificate")
	clientKey := flag.String("client-key", "", "path to the CA private key")
	maxCrashDuration := flag.Int("crash", 0, "maximum duration to wait before crashing")
	flag.Parse()

	if *maxCrashDuration != 0 {
		setupCrashRoutine(*maxCrashDuration)
	}

	log.Printf("Listening on [%v]...\n", *listenAddr)

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/healthz", handleHealthzRequest)

	var err error
	if *clientCert != "" && *clientKey != "" {
		err = http.ListenAndServeTLS(*listenAddr, *clientCert, *clientKey, nil)
	} else {
		err = http.ListenAndServe(*listenAddr, nil)
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

	log.Printf("Crashing in [%v] seconds", crashDuration)
	go func() {
		time.Sleep(time.Duration(crashDuration) * time.Second)
		log.Fatal("Crashing...")
	}()
}
