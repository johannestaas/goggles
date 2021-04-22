package connection

import (
    "net"
    "bufio"
    "log"
)


func HandleConnection(conn net.Conn) {
    buffer, err := bufio.NewReader(conn).ReadBytes('\n')

    if err != nil {
        log.Println("Client disconnected: ", err.Error())
        conn.Close()
        return
    }

    log.Println("client message: ", string(buffer[:len(buffer) - 1]))

    conn.Write(buffer)

    HandleConnection(conn)
}
