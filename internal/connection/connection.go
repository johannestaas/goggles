package connection

import (
	"bufio"
	"errors"
	"internal/kvstore"
	"log"
	"net"
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
		log.Printf("for some reason, the kvstore is nil\n")
		return "", errors.New("nil kvstore")
	}
	log.Printf("handling command %s args %#v\n", command, args)
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
		if err != nil {
			log.Printf("Client disconnected: %s\n", err.Error())
			conn.Close()
			return
		}
		command = command[:len(command)-1]

		cmd_len, err := commandArgLen(command)
		if err != nil {
			// yeah whateva i'll handle this later
			log.Printf("error: %s\n", err.Error())
			log.Printf("client sent bullshit, killing their conn: %s\n", command)
			conn.Close()
			return
		}

		var args []string
		for i := 0; i < cmd_len; i++ {
			var delim byte = ' '
			if i == cmd_len-1 {
				delim = '\n'
			}
			arg, err := rdr.ReadString(delim)
			if err != nil {
				log.Printf("Client disconnected: %s\n", err.Error())
				conn.Close()
				return
			}
			log.Printf("parsed arg %s\n", arg)
			args = append(args, arg[:len(arg)-1])
			log.Printf("args is %#v\n", args)
		}

		if command == "db" {
			db, ok := stores[args[0]]
			if ok {
				store = db
			} else {
				log.Printf("creating db %s\n", args[0])
				store = kvstore.New(args[0])
				stores[args[0]] = store
			}
			conn.Write([]byte("db set to " + args[0] + "\n"))
		} else {
			result, err := handleCommand(store, command, args)
			if err != nil {
				log.Printf("client sent more bullshit: %s\n", err.Error())
				conn.Close()
				return
			}
			conn.Write([]byte(result + "\n"))
		}
	}
}
