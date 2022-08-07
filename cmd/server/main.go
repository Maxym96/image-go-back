package main

import (
	"context"
	"errors"
	"fmt"
	"image-go-back/config"
	"image-go-back/config/container"
	"image-go-back/internal/infra/database"
	"image-go-back/internal/infra/http"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func main() {
	exitCode := 0
	ctx, cancel := context.WithCancel(context.Background())

	// Recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("The system panicked!: %v\n", r)
			fmt.Printf("Stack trace form panic: %s\n", string(debug.Stack()))
			exitCode = 1
		}
		os.Exit(exitCode)
	}()

	// Signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Printf("Received signal '%s', stopping... \n", sig.String())
		cancel()
		fmt.Printf("Sent cancel to all threads...")
	}()

	var conf = config.GetConfiguration()

	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	_, err = os.Stat(conf.FileStorageLocation)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(conf.FileStorageLocation, os.ModePerm)
		if err != nil {
			log.Fatalf("storage folder can not be created %s", err)
		}
	} else if err != nil {
		log.Fatalf("storage folder is not available %s", err)
	}

	cont := container.New(conf)

	// HTTP Server
	err = http.Server(
		ctx,
		http.Router(cont),
	)

	if err != nil {
		fmt.Printf("http server error: %s", err)
		exitCode = 2
		return
	}
}
