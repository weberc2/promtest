package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	pz "github.com/weberc2/httpeasy"
	"github.com/weberc2/httpeasy/promrouter"
)

func main() {
	var (
		promAddr = getEnv("PROM_ADDR", ":9091")
		appAddr  = getEnv("ADDR", ":8080")
	)

	log.Printf("Starting metrics server at %s", promAddr)
	go func() {
		var mux http.ServeMux
		mux.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(promAddr, &mux); err != nil {
			log.Fatalf("Starting metrics server: %v", err)
		}
	}()

	log.Printf("Starting app server at %s", appAddr)
	if err := http.ListenAndServe(
		appAddr,
		promrouter.NewWithDefaults().Register(
			// like pz.JSONLog(), but this doesn't pretty-print
			func(v interface{}) {
				data, err := json.Marshal(v)
				if err != nil {
					log.Printf("ERROR marshaling %# v", v)
					return
				}
				fmt.Fprintf(os.Stderr, "%s\n", data)
			},
			pz.Route{Path: "/health", Method: "GET", Handler: healthCheck},
		),
	); err != nil {
		log.Fatalf("Starting app server: %v", err)
	}
}

func healthCheck(pz.Request) pz.Response { return pz.Ok(pz.String("OK")) }

func getEnv(env, defaultValue string) string {
	if x := os.Getenv(env); x != "" {
		return x
	}
	return defaultValue
}
