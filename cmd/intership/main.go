package main

import "github.com/PIRSON21/mediasoft-intership2025/internal/server"

const version = "v1.0"

func main() {
	server.CreateServer(version)
}
