package main

import (
	"fmt"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/client"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
)

var (
	// BuildVersion is a build version of client application.
	BuildVersion = "N/A"
	// BuildDate is a build date of client application.
	BuildDate = "N/A"
	// BuildCommit is a build commit of client application.
	BuildCommit = "N/A"
)

func main() {
	// print server build version
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	if err := client.RunClient(cfg); err != nil {
		panic(err)
	}
}
