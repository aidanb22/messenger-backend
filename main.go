package main

import "github.com/ablancas22/messenger-backend/cmd"

func main() {
	var app cmd.App
	err := app.Initialize()
	if err != nil {
		panic(err)
	}
	app.Run()
}
