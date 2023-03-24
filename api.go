package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiServer struct {
	listenAddr string
	certFile   string
	keyFile    string
}

func NewApiServer(config Config) ApiServer {
	return ApiServer{
		listenAddr: config.listenAddr,
		certFile:   config.certFile,
		keyFile:    config.keyFile,
	}
}

func (s *ApiServer) Run() {
	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		fmt.Println("Unable to load cert or key file")
		panic(err)
	}

	fmt.Println("Starting webhook server")
	http.HandleFunc("/validate", validate)
	server := http.Server{
		Addr: s.listenAddr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}

func validate(w http.ResponseWriter, r *http.Request) {
	if err := writeJson(w, http.StatusOK, map[string]string{
		"message": "hello there",
	}); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
}

func writeJson(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
