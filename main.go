package main

import (
	"fmt"
	"net"
	"strings"

	"go.uber.org/atomic"
)

var (
	red    = "\u001b[31m"
	yellow = "\u001b[33m"
	green  = "\u001b[32m"
	purple = "\u001b[35m"
	black  = "\u001b[30m"

	whiteBG = "\u001b[47m"

	clear = "\u001b[0m"
	bold  = "\u001b[1m"
)

func main() {
	listener, err := net.Listen("tcp", ":5050")

	if err != nil {
		panic(err)
	}

	for {
		client, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		host, err := net.Dial("tcp", "127.0.0.1:7656")

		if err != nil {
			panic(err)
		}

		defer host.Close()

		ID := atomic.NewString("undefined")

		go pipe(ID, fmt.Sprintf("%sclient -> host%s", purple, clear), client, host)
		go pipe(ID, fmt.Sprintf("%shost -> client%s", green, clear), host, client)
	}
}

func pipe(ID *atomic.String, title string, a net.Conn, b net.Conn) {
	buffer := make([]byte, 65535)

	for {
		length, err := a.Read(buffer)

		if err != nil {
			issue(err)
			return
		}

		cursor := buffer[:length]

		var formatted string

		for index, parameter := range strings.Split(string(cursor), " ") {
			split := strings.Split(parameter, "=")

			if len(split) == 2 {
				formatted += fmt.Sprintf("NEWLINETAB%sTAB%s%s%s", split[0], red, split[1], clear)

				if split[0] == "ID" {
					ID.Store(split[1])
				}
			} else {
				if index != 0 {
					formatted += " "
				}

				formatted += parameter
			}
		}

		formatted = strings.ReplaceAll(formatted, "\n", yellow+"\\\\n"+clear)
		formatted = strings.ReplaceAll(formatted, "\t", yellow+"\\\\t"+clear)
		formatted = strings.ReplaceAll(formatted, "\r", yellow+"\\\\r"+clear)

		formatted = strings.ReplaceAll(formatted, "NEWLINE", "\n")
		formatted = strings.ReplaceAll(formatted, "TAB", "\t")

		fmt.Printf("%s%s%s %s %s\n", bold, ID.Load(), clear, title, formatted)

		_, err = b.Write(cursor)

		if err != nil {
			issue(err)
			return
		}
	}
}

func issue(err error) {
	fmt.Printf("%s%sunexpected issue: %v%s\n", black, whiteBG, err, clear)
}
