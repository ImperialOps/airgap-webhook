package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiServer interface {
	Run()
}

type ApiServerCommon struct {
	config  *Config
	backend IBackend
}

type ApiServerHttp struct {
	ApiServerCommon
}

type ApiServerHttps struct {
	ApiServerCommon
}

func NewApiServer(c *Config) ApiServer {
	apiServer := ApiServerCommon{
		config:  c,
		backend: NewBackend(c.backend),
	}

	switch c.tls.enabled {
	case true:
		return &ApiServerHttps{
			ApiServerCommon: apiServer,
		}
	default:
		return &ApiServerHttp{
			ApiServerCommon: apiServer,
		}
	}
}

func (s *ApiServerHttps) Run() {
	cert, err := tls.LoadX509KeyPair(s.config.tls.certFile, s.config.tls.keyFile)
	if err != nil {
		log.Println("Unable to load cert or key file")
		panic(err)
	}

	log.Printf("listening on %s", s.config.listenAddr)
	http.HandleFunc("/healthz", newApiFunc(s.handleHealth))
	http.HandleFunc("/validate", newApiFunc(s.handleValidate))
	server := http.Server{
		Addr: s.config.listenAddr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}

func (s *ApiServerHttp) Run() {
	log.Printf("listening on %s", s.config.listenAddr)
	http.HandleFunc("/healthz", newApiFunc(s.handleHealth))
	http.HandleFunc("/validate", newApiFunc(s.handleValidate))
	server := http.Server{
		Addr: s.config.listenAddr,
	}

	if err := server.ListenAndServe(); err != nil {
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

func (s *ApiServerCommon) handleValidate(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handlePostValidate(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, fmt.Sprintf("%s method not allowed", r.Method))
	}
}

func (s *ApiServerCommon) handlePostValidate(w http.ResponseWriter, r *http.Request) error {
	// Validate that the incoming content type is correct.
	if r.Header.Get("Content-Type") != "application/json" {
		return NewApiError(http.StatusBadRequest, "expected application/json content-type")
	}

	// Get the body data, which will be the AdmissionReview
	// content for the request.
	var body []byte
	if r.Body != nil {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return NewApiError(http.StatusBadRequest, "expected a request body")
		}
		body = requestData
	}

	response, err := handleAdmissionReview(body)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, response)
}

func (s *ApiServerCommon) handleHealth(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetHealth(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, fmt.Sprintf("%s method not allowed", r.Method))
	}
}

func (s *ApiServerCommon) handleGetHealth(w http.ResponseWriter, r *http.Request) error {
	if s.isServerHealthy() {
		return writeJson(w, http.StatusOK, "")
	} else {
		return writeJson(w, http.StatusInternalServerError, "")
	}
}

func (s *ApiServerCommon) isServerHealthy() bool {
	// TODO check DB connection
	return true
}

func writeJson(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
