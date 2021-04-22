package main

import (
    "fmt"
    "os"
    "net"
    "internal/connection"
)


const (
    connHost = "localhost"
    connPort = "7712"
    connType = "tcp"
)


func main() {
    var connStr = connHost + ":" + connPort
    fmt.Println("Starting " + connStr)
    sock, err := net.Listen(connType, connStr)
    if err != nil {
        fmt.Println("error listening: ", err.Error())
        os.Exit(1)
    }
    defer sock.Close()

    for {
        conn, err := sock.Accept()
        if err != nil {
            fmt.Println("Error connecting: ", err.Error())
            return
        }
        fmt.Println("client connected")
        fmt.Println("Client " + conn.RemoteAddr().String() + " connected")

        go connection.HandleConnection(conn)
    }
}
