package main

import (
	"flag"
	// "fmt"
	"lan-cloud/internal/config"
	"lan-cloud/internal/filesystem"
	"lan-cloud/internal/metadata"
	"lan-cloud/internal/server"
	// "time"
)

var (
	cliStoragePath string
	resetPassword string
)

func init() {
	flag.StringVar(&cliStoragePath, "storage", "", "Path to storage directory")
	flag.StringVar(&resetPassword, "reset", "", "Reset password")
}

func main() {
	flag.Parse()
	// start := time.Now()
	if err := config.LoadStorageConfig(cliStoragePath); err != nil { panic(err) }
	if err := config.LoadServerPassConfig(resetPassword); err != nil { panic(err) }
	if err := metadata.Load(); err != nil { panic(err) }
	if err := filesystem.SyncMetadata(); err != nil { panic(err) }
	// elapsed := time.Since(start)
	// fmt.Printf("Process took %s\n", elapsed)

	server.Start()

	// fmt.Println("Press enter to exit...")
	// fmt.Scanln()
}