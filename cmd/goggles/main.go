package main

import (
    "fmt"
    "flag"
    "log"
    "os"
    "net"
    "internal/connection"
    "internal/kvstore"
)


const (
    defaultHost = "localhost"
    defaultPort = 7712
)


func main() {
    host := flag.String("host", defaultHost, "bind to this host")
    port := flag.Int("port", defaultPort, "the port to serve on")
    flag.Parse()
    connStr := fmt.Sprintf("%s:%d", *host, *port)
    stores := make(map[string]*kvstore.KVStore)
    stores["test"] = kvstore.New("test")
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
            return
        }
        log.Printf("client %s connected\n", conn.RemoteAddr().String())

        go connection.HandleConnection(conn, stores)
    }
}
