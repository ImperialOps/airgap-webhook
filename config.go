package main

type Config struct {
	listenAddr string
	certFile   string
	keyFile    string
}

func NewConfig() Config {
	return Config{
		listenAddr: ":8000",
		certFile:   "",
		keyFile:    "",
	}
}
