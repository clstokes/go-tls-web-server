package main

import (
	"flag"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/circonus"
)

const (
	ResponseOk = "{\"status\":\"ok\"}"
)

type Server struct {
	listenAddr       string
	clientCert       string
	clientKey        string
	maxCrashDuration int

	circonusConfig circonus.Config
	metrics        metrics.FanoutSink
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	server := &Server{}
	server.parseArgs()

	server.setupMetricsSink()
	server.setupCrashRoutine()

	server.serve()
	return 0
}

func (c *Server) serve() {
	http.HandleFunc("/", c.handleRequest)
	http.HandleFunc("/healthz", c.handleHealthzRequest)

	log.Printf("Listening on [%v]...\n", c.listenAddr)

	var err error
	if c.clientCert != "" && c.clientKey != "" {
		err = http.ListenAndServeTLS(c.listenAddr, c.clientCert, c.clientKey, nil)
	} else {
		err = http.ListenAndServe(c.listenAddr, nil)
	}

	if err != nil {
		log.Fatal("Error starting HTTP listener: ", err)
	}
}

func (c *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request [%v]...\n", html.EscapeString(r.URL.Path))
	w.Write([]byte(ResponseOk))
	c.metrics.IncrCounter([]string{html.EscapeString(r.URL.Path)}, 1)
	return
}

func (c *Server) handleHealthzRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(ResponseOk))
	return
}

func (c *Server) setupMetricsSink() {
	c.metrics = metrics.FanoutSink{}
	if c.circonusConfig.CheckManager.API.TokenKey == "" {
		log.Printf("Circonus API Token not provided. Using in-memory metrics sink.")
		c.metrics = append(c.metrics, metrics.NewInmemSink(10*time.Second, time.Minute))
		return
	}

	log.Printf("Configuring Circonus metrics Sink for [%v]...", c.circonusConfig.CheckManager.API.URL)
	sink, err := circonus.NewCirconusSink(&c.circonusConfig)
	if err != nil {
		log.Fatalf("Error setting up Circonus metrics Sink for [%v]: ", err)
	}
	sink.Start()

	c.metrics = append(c.metrics, sink)
}

func (c *Server) parseArgs() {
	// main args
	flag.StringVar(&c.listenAddr, "listen-address", ":8080", "address to listen on")
	flag.StringVar(&c.clientCert, "client-cert", "", "path to the CA certificate (default \"\")")
	flag.StringVar(&c.clientKey, "client-key", "", "path to the CA private key (default \"\")")
	flag.IntVar(&c.maxCrashDuration, "crash", 0, "maximum duration to wait before crashing (default \"0\" - e.g. don't crash on purpose)")

	// circonus args
	flag.StringVar(&c.circonusConfig.CheckManager.API.TokenKey, "circonus-api-token", "", "circonus - API token. If provided, metric management is enabled. (default \"\")")
	flag.StringVar(&c.circonusConfig.CheckManager.API.URL, "circonus-api-url", "", "circonus - base URL for the API (default \"\")")
	flag.StringVar(&c.circonusConfig.CheckManager.API.TokenApp, "circonus-api-app", "", "circonus - app name associated with the API token (default \"\")")

	flag.Parse()
}

func (c *Server) setupCrashRoutine() {
	if c.maxCrashDuration == 0 {
		return
	}

	rand.Seed(time.Now().Unix())
	crashDuration := rand.Intn(c.maxCrashDuration)

	log.Printf("Crashing in [%v] seconds...", crashDuration)
	go func() {
		time.Sleep(time.Duration(crashDuration) * time.Second)
		log.Fatal("Crashing NOW...")
	}()
}
