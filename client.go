package main

type IBackend interface {
	Send([]Image) error
}

type HttpClient struct {
	config ConfigBackend
}

func NewBackend(config ConfigBackend) IBackend {
	switch config.protocol {
	case "http":
		return NewHttpClient(config)
	default:
		return nil
	}
}

func NewHttpClient(config ConfigBackend) *HttpClient {
	return &HttpClient{
		config: config,
	}
}

func (c *HttpClient) Send([]Image) error {
	return nil
}
