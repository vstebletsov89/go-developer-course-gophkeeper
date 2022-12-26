package main

import (
	"fmt"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/server"
)

var (
	// BuildVersion is a build version of server application.
	BuildVersion = "N/A"
	// BuildDate is a build date of server application.
	BuildDate = "N/A"
	// BuildCommit is a build commit of server application.
	BuildCommit = "N/A"
)

func main() {
	// print server build info
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	if err := server.RunServer(cfg); err != nil {
		panic(err)
	}
}
