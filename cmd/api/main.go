package main

import (
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/server"
)

//	@title			Go Newsletter platform API
//	@version		1.0
//	@description	This is a sample RESTful API for a Go Newsletter platform

//	@contact.name	Mohammed Shakleya
//	@contact.url	https://www.linkedin.com/in/mohammed-shakleya/

// @host		localhost:8080
// @basePath	/v1
func main() {

	// Get a new server instance.
	s := server.New()

	// Initialize server instance.
	s.Init()

	// Run
	s.Run()
}