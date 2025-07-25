package main

import (
	"flag"
	"fmt"
	"lan-cloud/internal/config"
	"lan-cloud/internal/filesystem"
	"lan-cloud/internal/metadata"
	"lan-cloud/internal/server"
	"time"
)

var (
	cliStoragePath string
	resetPassword string
	skipSync bool
)

func init() {
	flag.StringVar(&cliStoragePath, "storage", "", "Path to storage directory")
	flag.StringVar(&resetPassword, "reset", "", "Reset password")
	flag.BoolVar(&skipSync, "skip-sync", false, "Skip sync")
}

func main() {
	flag.Parse()
	
	if err := config.LoadStorageConfig(cliStoragePath); err != nil { panic(err) }
	if err := config.LoadServerPassConfig(resetPassword); err != nil { panic(err) }
	
	fmt.Println("Loading Metadata...")
	metadataProcess := time.Now()
	if err := metadata.Load(); err != nil { panic(err) }
	fmt.Printf("Metadata Process took %s\n", time.Since(metadataProcess))
	
	if !skipSync {
		fmt.Println("Syncing Metadata...")
		syncProcess := time.Now()
		if err := filesystem.SyncMetadata(); err != nil { panic(err) } 
		fmt.Printf("Sync Process took %s\n", time.Since(syncProcess))
	}

	server.Start()
}