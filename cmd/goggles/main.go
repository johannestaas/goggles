package main

import (
	"flag"
	"fmt"
	"internal/connection"
	"internal/database"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultHost    = "localhost"
	defaultPort    = 7712
	defaultDataDir = "/var/db/goggles"
)

func checkDataDirOrExit(dataDir *string) {
	if _, err := os.Stat(*dataDir); os.IsNotExist(err) {
		log.Printf("data dir %s does not exist, exitting.\n", *dataDir)
		os.Exit(1)
	}
}

func main() {
	// Handle ctrl-c
	ctrlC := make(chan os.Signal)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctrlC
		log.Println("caught ctrl-c, quitting...")
		os.Exit(0)
	}()

	// cli args
	host := flag.String("host", defaultHost, "bind to this host")
	port := flag.Int("port", defaultPort, "the port to serve on")
	dataDir := flag.String("datadir", defaultDataDir, "the directory to persist data")
	flag.Parse()

	checkDataDirOrExit(dataDir)

	connStr := fmt.Sprintf("%s:%d", *host, *port)

	// init db
	db := database.New(dataDir)

	// start up server
	log.Printf("Starting %s\n", connStr)
	sock, err := net.Listen("tcp", connStr)
	if err != nil {
		log.Printf("error listening: %s\n", err.Error())
		os.Exit(1)
	}
	defer sock.Close()

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Printf("Error connecting: %s\n", err.Error())
			continue
		}
		log.Printf("client %s connected\n", conn.RemoteAddr().String())
		go connection.HandleConnection(conn, db)
	}
}
