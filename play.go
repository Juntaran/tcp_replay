package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var bindAddress string
var sourceName string

func init() {
	flag.StringVar(&bindAddress, "bind", ":2000", "The listen address")
	flag.StringVar(&sourceName, "source", "-", "The recorded source file timestamp-nano ff ff ff ff etc\\n")
}

func main() {
	flag.Parse()
	err := do()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
		return
	}

}

func broadcast(msg []byte) {
	for id, c := range listeners {
		go func(id string, c net.Conn) {
			_, err := c.Write(msg)
			if err != nil {
				log.Println(err.Error())
				delete(listeners, id)
			}
		}(id, c)
	}
}

var listeners = make(map[string]net.Conn)
var _lastID int
var _lastIDMutex = sync.Mutex{}

func nextID() string {
	_lastIDMutex.Lock()
	defer _lastIDMutex.Unlock()
	_lastID += 1
	return fmt.Sprintf("%d", _lastID)
}

func addConnections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		listeners[nextID()] = conn
	}
}

func do() error {
	var source io.Reader
	if sourceName == "-" {
		source = os.Stdin
	} else {
		file, err := os.Open(sourceName)
		if err != nil {
			return err
		}
		defer file.Close()
		source = file
	}

	fmt.Printf("Reading from %s\n", sourceName)

	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return err
	}
	go addConnections(listener)
	for {
		readFile(source)
	}
	return nil
}

func readFile(source io.Reader) {
	scanner := bufio.NewScanner(source)
	var startFileTimestamp int64
	var startTimestamp int64 = time.Now().UnixNano()
	for scanner.Scan() {
		s := scanner.Text()
		parts := strings.Split(s, " ")
		ts, _ := strconv.ParseInt(parts[0], 10, 64)
		if startFileTimestamp > 0 {
			currentTick := time.Now().UnixNano() - startTimestamp
			fileTick := ts - startFileTimestamp
			delay := fileTick - currentTick
			if delay > 0 {
				time.Sleep(time.Duration(delay) * time.Nanosecond)
			}
		} else {
			startFileTimestamp = ts
		}

		byteArray := make([]byte, len(parts)-1, len(parts)-1)
		for i, b := range parts[1:] {
			iVal, _ := strconv.ParseInt(b, 16, 16)
			byteArray[i] = byte(iVal)
		}
		fmt.Println(time.Unix(ts/1000000000, 0).Format(time.RFC3339Nano))
		fmt.Printf("%d % x\n", ts, byteArray)
		broadcast(byteArray)
	}
}
