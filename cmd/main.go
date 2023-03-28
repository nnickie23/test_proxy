package main

import "github.com/nnickie23/test_proxy/internal/app"

//swagger

// @title TestTask
// @version 1.0
// @description Requestor API application

// host localhost:8000
// @host localhost:8000
// @schemes http https
// @Basepath /
func main() {
	app.Run()
}
