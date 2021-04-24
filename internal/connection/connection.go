package connection

import (
    "net"
    "bufio"
    "log"
    "errors"
    "internal/kvstore"
)


func commandArgLen(command string) (int, error) {
    switch command {
        case "db":
            return 1, nil
        case "get":
            return 1, nil
        case "set":
            return 2, nil
        default:
            return 0, errors.New("unknown command " + command)
    }
}


func handleCommand(store *kvstore.KVStore, command string, args []string) (string, error) {
    if store == nil {
        log.Println("for some reason, the kvstore is nil")
        return "", errors.New("nil kvstore")
    }
    log.Println("handling command %s args %v", command, args)
    switch command {
        case "get":
            result := store.Get(args[0])
            return result, nil
        case "set":
            store.Set(args[0], args[1])
            return "done", nil
        default:
            return "", errors.New("unknown command " + command)
    }
}


func HandleConnection(conn net.Conn, stores map[string]*kvstore.KVStore) {
    var store *kvstore.KVStore = nil
    for {
        rdr := bufio.NewReader(conn)
        command, err := rdr.ReadString(' ')
        command = command[:len(command) - 1]
        if err != nil {
            log.Println("Client disconnected: ", err.Error())
            conn.Close()
            return
        }

        cmd_len, err := commandArgLen(command)
        if err != nil {
            // yeah whateva i'll handle this later
            log.Println("error: %s", err.Error())
            log.Println("client sent bullshit, killing their conn: %s", command)
            conn.Close()
            return
        }

        var args []string
        for i := 0; i < cmd_len; i++ {
            var delim byte = ' '
            if i == cmd_len - 1 {
                delim = '\n'
            }
            arg, err := rdr.ReadString(delim)
            if err != nil {
                log.Println("Client disconnected: %s", err.Error())
                conn.Close()
                return
            }
            log.Println("parsed arg %s", arg)
            args = append(args, arg[:len(arg) - 1])
            log.Println("args is %v", args)
        }

        if command == "db" {
            db, ok := stores[args[0]]
            if ok {
                store = db
            } else {
                log.Println("creating db %s", args[0])
                store = kvstore.New(args[0])
                stores[args[0]] = store
            }
            conn.Write([]byte("db set to " + args[0] + "\n"))
        } else {
            result, err := handleCommand(store, command, args)
            if err != nil {
                log.Println("client sent more bullshit: %s", err.Error())
                conn.Close()
                return
            }
            conn.Write([]byte(result + "\n"))
        }
    }
}
