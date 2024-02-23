package main

import (
	"log"

	app "github.com/ingeniousambivert/fiber-bootstrapped/src/app"
	core "github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func main() {
	server := core.Build()
	app.Init(server)

	err := server.Boot()
	if err != nil {
		log.Fatalf("server:error: failed to start server %s", err)
	}
}
