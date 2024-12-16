package main

import (
	"github.com/OinkiePie/calc/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
	// app.Run()
}