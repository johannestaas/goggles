package connection

import (
	"bufio"
	"errors"
	"internal/database"
	"internal/kvstore"
	"log"
	"net"
	"strconv"
	"time"
)

const timeout = 5 * time.Second

func commandArgLen(command string) (int, error) {
	switch command {
	case "db":
		return 1, nil
	case "get":
		return 1, nil
	case "set":
		return 3, nil
	default:
		return 0, errors.New("unknown command " + command)
	}
}

func handleGet(store *kvstore.KVStore, key string) (string, error) {
	result := store.Get(key)
	return result, nil
}

func handleSet(store *kvstore.KVStore, durStr string, key string, val string) (string, error) {
	duration, err := strconv.ParseUint(durStr, 10, 64)
	if err != nil {
		return "", errors.New("duration was not an int: " + durStr)
	}
	store.Set(key, val, time.Duration(duration)*time.Second)
	return "", nil
}

func handleCommand(store *kvstore.KVStore, command string, args []string) (string, error) {
	if store == nil {
		log.Printf("for some reason, the kvstore is nil\n")
		return "", errors.New("nil kvstore")
	}
	log.Printf("handling command %s args %#v\n", command, args)
	switch command {
	case "get":
		return handleGet(store, args[0])
	case "set":
		return handleSet(store, args[0], args[1], args[2])
	default:
		return "", errors.New("unknown command " + command)
	}
}

func makeArgs(rdr *bufio.Reader, cmdLen int) ([]string, error) {
	var args []string
	for i := 0; i < cmdLen; i++ {
		var delim byte = ' '
		if i == cmdLen-1 {
			delim = '\n'
		}
		arg, err := rdr.ReadString(delim)
		if err != nil {
			log.Printf("Client disconnected: %s\n", err.Error())
			return args, errors.New("client disconnected")
		}
		log.Printf("parsed arg %s\n", arg)
		args = append(args, arg[:len(arg)-1])
		log.Printf("args is %#v\n", args)
	}
	return args, nil
}

func HandleConnection(conn net.Conn, db *database.Database) {
	// stores := *map[string]*kvstore.KVStore
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

		cmdLen, err := commandArgLen(command)
		if err != nil {
			// yeah whateva i'll handle this later
			log.Printf("error: %s\n", err.Error())
			log.Printf("command not recognized: %s\n", command)
			conn.Write([]byte("error: bad command\n"))
			log.Println("resetting connection bufio reader")
			rdr.Reset(conn)
			continue
		}

		args, err := makeArgs(rdr, cmdLen)
		if err != nil {
			log.Println("client disconnected")
			conn.Close()
			return
		}

		if command == "db" {
			store = db.GetOrCreateStore(&args[0])
			conn.Write([]byte("\n"))
		} else {
			result, err := handleCommand(store, command, args)
			if err != nil {
				log.Println("bad command, resetting connection bufio reader")
				conn.Write([]byte("error: bad command string\n"))
				rdr.Reset(conn)
				continue
			}
			conn.Write([]byte(result + "\n"))
		}
	}
}
