package main

import (
	app "github.com/savisitor15/db-pokedex-go/internal/app"
)

func main() {
	app.InitializeState()
	app.CommandLoop()
}
