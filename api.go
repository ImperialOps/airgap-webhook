package main

import (
	"crypto/tls"
	"encoding/json"
    "errors"
    "log"
	"fmt"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiServer struct {
	listenAddr string
	certFile   string
	keyFile    string
}

func NewApiServer(config Config) *ApiServer {
	return &ApiServer{
		listenAddr: config.listenAddr,
		certFile:   config.certFile,
		keyFile:    config.keyFile,
	}
}

func (s *ApiServer) Run() {
	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		log.Println("Unable to load cert or key file")
		panic(err)
	}

	log.Println("Starting webhook server")
	http.HandleFunc("/validate", newApiFunc(s.handleValidate))
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

func newApiFunc(f apiFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := f(w, r); err != nil {
            var apiError *ApiError
            if errors.As(err, &apiError) {
                writeJson(w, apiError.Code(), apiError.Error())
                return
            }
            writeJson(w, http.StatusInternalServerError, err.Error())
        }
    }
}

func (s *ApiServer) handleValidate(w http.ResponseWriter, r *http.Request) error {
    switch r.Method {
    case "GET":
        return s.handleGetValidate(w, r)
    case "POST":
        return s.handlePostValidate(w, r)
    default:
        return NewApiError(http.StatusMethodNotAllowed, fmt.Sprintf("%s method not allowed", r.Method))
    }
}

func (s *ApiServer) handleGetValidate(w http.ResponseWriter, r *http.Request) error {
    return writeJson(w, http.StatusOK, "okkkkk")
}

func (s *ApiServer) handlePostValidate(w http.ResponseWriter, r *http.Request) error {
    return nil
}

func writeJson(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
