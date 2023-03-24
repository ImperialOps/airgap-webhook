package main

func main() {
	config := NewConfig()

	server := NewApiServer(config)
	server.Run()
}
