package server

import (
	"crypto/tls"
	"log"
	"net/http"
)

func Run() error {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.Println("Listening on :443")
	return srv.ListenAndServeTLS("/certs/tls.crt", "/certs/tls.key")
}
