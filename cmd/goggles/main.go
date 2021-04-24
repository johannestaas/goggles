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
    log.Println("Starting %s", connStr)
    sock, err := net.Listen(connType, connStr)
    if err != nil {
        log.Println("error listening: %s", err.Error())
        os.Exit(1)
    }
    defer sock.Close()

    for {
        conn, err := sock.Accept()
        if err != nil {
            log.Println("Error connecting: %s", err.Error())
            return
        }
        log.Println("client %s connected", conn.RemoteAddr().String())

        go connection.HandleConnection(conn, stores)
    }
}
