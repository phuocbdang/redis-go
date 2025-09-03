package main

import "redisgo/internal/server"

func main() {
	server.RunIOMultiplexingServer()
}
