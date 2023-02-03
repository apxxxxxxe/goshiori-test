package main

/*
REQUEST SHIORIPROXY/1.0
GET SHIORI/3.0
Charset: Shift_JIS
Sender: ikagaka
ID: version

LOAD SHIORIPROXY/1.0
C:\SSP\ghost\ikaga\ghost\master\

*/

import (
	"bufio"
	"fmt"
	"os"

	shiori "github.com/Narazaka/shiorigo"
	"github.com/apxxxxxxe/gohst/dictionary"
)

const name = "GO"
const version = "0.1.0"
const logFile = "err.log"

const (
	statusOK        = 200
	statusNoContent = 204
)

var (
	Dictionary *dictionary.Dictionary
)

func main() {
	var scanner = bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		switch scanner.Text() {
		case "LOAD SHIORIPROXY/1.0":
			load(scanner)
		case "REQUEST SHIORIPROXY/1.0":
			request(scanner)
		case "UNLOAD SHIORIPROXY/1.0":
			unload()
		}
	}
}

func load(scanner *bufio.Scanner) string {
	scanner.Scan()
	dirPath := scanner.Text()
	fmt.Print("1\r\n")
	return dirPath
}

func request(scanner *bufio.Scanner) {
	requestStr := ""
	for {
		scanner.Scan()
		line := scanner.Text()
		requestStr += line + "\r\n"
		if scanner.Text() == "" {
			break
		}
	}
	r, err := shiori.ParseRequest(requestStr)
	if err != nil {
		os.WriteFile(logFile, []byte(err.Error()), 0644)
	}
	processRequest(r)
}

func processRequest(req shiori.Request) {
	var value string
	status := statusOK

	switch req.Headers["ID"] {
	case "version":
		value = version
	case "OnFirstBoot", "OnBoot":
		value = "\\0HogeFuga\\e"
	}

	res := shiori.Response{
		Protocol: shiori.SHIORI,
		Version:  "3.0",
		Code:     status,
		Headers: shiori.ResponseHeaders{
			"Charset": "Shift_JIS",
			"Sender":  name,
			"Value":   value,
		},
	}

	fmt.Print(res)
}

func unload() {
	fmt.Print("1\r\n")
	os.Exit(0)
}
