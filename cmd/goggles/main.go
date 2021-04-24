package main

import (
    "log"
    "os"
    "net"
    "internal/connection"
    "internal/kvstore"
)


const (
    connHost = "localhost"
    connPort = "7712"
    connType = "tcp"
)


func main() {
    connStr := connHost + ":" + connPort
    stores := make(map[string]*kvstore.KVStore)
    stores["test"] = kvstore.New("test")
    log.Printf("Starting %s\n", connStr)
    sock, err := net.Listen(connType, connStr)
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
