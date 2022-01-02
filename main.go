package main

import (
	"fmt"
	"net"
	"strings"
)

var (
	red    = "\u001b[31m"
	yellow = "\u001b[33m"
	green  = "\u001b[32m"
	purple = "\u001b[35m"
	clear  = "\u001b[0m"
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

		go pipe(fmt.Sprintf("%sclient -> host%s", purple, clear), client, host)
		go pipe(fmt.Sprintf("%shost -> client%s", green, clear), host, client)
	}
}

func pipe(title string, a net.Conn, b net.Conn) {
	buffer := make([]byte, 65535)

	for {
		length, err := a.Read(buffer)

		if err != nil {
			fmt.Println(err)
			return
		}

		cursor := buffer[:length]

		var formatted string

		for index, parameter := range strings.Split(string(cursor), " ") {
			split := strings.Split(parameter, "=")

			if len(split) == 2 {
				formatted += fmt.Sprintf("NEWLINETAB%sTAB%s%s%s", split[0], red, split[1], clear)
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

		fmt.Printf("%s: %s\n", title, formatted)

		_, err = b.Write(cursor)

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
