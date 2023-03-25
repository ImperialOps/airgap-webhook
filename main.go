package main

func main() {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	server := NewApiServer(config)
	server.Run()
}
