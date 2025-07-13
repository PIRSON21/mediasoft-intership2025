package main

import "github.com/PIRSON21/mediasoft-intership2025/internal/server"

const version = "v0.0.5"

func main() {
	server.CreateServer(version)
}
