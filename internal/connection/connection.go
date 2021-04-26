package connection

import (
	"bufio"
	"errors"
	"fmt"
	"internal/database"
	"internal/kvstore"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const timeout = 5 * time.Second

func commandArgLen(command *string) (int, error) {
	switch *command {
	case "db":
		return 1, nil
	case "get":
		return 1, nil
	case "set":
		return 3, nil
	case "drop":
		return 0, nil
	default:
		log.Printf("client error: got the command %v\n", *command)
		return 0, errors.New("unknown command " + *command)
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

func handleDrop(db *database.Database, store *kvstore.KVStore) (string, error) {
	if store == nil {
		log.Printf("dropping kvstore but it's nil\n")
		return "", errors.New("nil kvstore")
	}
	db.DropStore(store)
	return "", nil
}

func handleCommand(db *database.Database, store *kvstore.KVStore, command string, args []string) (string, error) {
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
	case "drop":
		return handleDrop(db, store)
	default:
		return "", errors.New("unknown command " + command)
	}
}

func makeArgs(argStr string, cmdLen int) ([]string, error) {
	args := strings.SplitAfterN(argStr, " ", cmdLen)
	if len(args) == 1 && args[0] == "" {
		// Split will return an array of size 1 with "" if it's empty.
		args = []string{}
	}
	if len(args) != cmdLen {
		errMsg := fmt.Sprintf("expected %d arguments, got %s", cmdLen, argStr)
		return args, errors.New(errMsg)
	}
	for i, arg := range args {
		// Trim the leading " " or "\n"
		size := len(arg)
		args[i] = arg[:size-1]
	}
	return args, nil
}

func splitStatement(statement *string) (string, string) {
	split := strings.SplitAfterN(*statement, " ", 2)
	command := strings.TrimSpace(split[0])
	if len(split) == 1 {
		return command, ""
	}
	argStr := split[1]
	return command, argStr
}

func HandleConnection(conn net.Conn, db *database.Database) {
	var store *kvstore.KVStore = nil
	for {
		rdr := bufio.NewReader(conn)
		statement, err := rdr.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %s\n", err.Error())
			conn.Close()
			return
		}

		command, argStr := splitStatement(&statement)

		cmdLen, err := commandArgLen(&command)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
			log.Printf("command not recognized: %s\n", command)
			conn.Write([]byte("error: bad command\n"))
			log.Println("resetting connection bufio reader")
			rdr.Reset(conn)
			continue
		}

		args, err := makeArgs(argStr, cmdLen)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
			log.Println("client disconnected")
			conn.Close()
			return
		}

		if command == "db" {
			store = db.GetOrCreateStore(&args[0])
			conn.Write([]byte("\n"))
		} else {
			result, err := handleCommand(db, store, command, args)
			if err != nil {
				log.Println("bad command, resetting connection bufio reader")
				conn.Write([]byte("error: bad command string\n"))
				rdr.Reset(conn)
				continue
			}
			if command == "drop" {
				store = nil
			}
			conn.Write([]byte(result + "\n"))
		}
	}
}
