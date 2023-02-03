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
	"log"
	"os"
	"path/filepath"

	shiori "github.com/Narazaka/shiorigo"
	"github.com/apxxxxxxe/goshiori-test/dictionary"
)

var (
	Dictionary *dictionary.Dictionary
	logFile    *os.File
)

func main() {
	var scanner = bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		switch scanner.Text() {
		case "LOAD SHIORIPROXY/1.0":
			if err := load(scanner); err != nil {
				log.Printf("[error] load: %s\n", err)
			}
			log.Println("[info] load()")
			fmt.Print("1\r\n")

		case "REQUEST SHIORIPROXY/1.0":
			log.Println("[info] request()")
			res, err := request(scanner)
			if err != nil {
				log.Printf("[error] request: %s\n", err)
			}
			log.Println(res)
			fmt.Print(res)

		case "UNLOAD SHIORIPROXY/1.0":
			log.Println("[info] unload()")
			if err := unload(); err != nil {
				log.Printf("[error] unload: %s\n", err)
			}
			fmt.Print("1\r\n")
			os.Exit(0)
		}
	}
}

func load(scanner *bufio.Scanner) error {
	exec, err := os.Executable()
	if err != nil {
		return err
	}
	dicDir := filepath.Dir(exec)

	// 辞書と変数をDictionaryに読み込む
	Dictionary, err = dictionary.New()
	if err != nil {
		return err
	}

	//ログ書き込み設定
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
	path := filepath.Join(dicDir, "shiori.log")
	logFile, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err.Error())
	} else if err == nil && logFile != nil {
		log.SetOutput(logFile)
	}

	return nil
}

func request(scanner *bufio.Scanner) (shiori.Response, error) {
	requestStr := ""
	for {
		scanner.Scan()
		line := scanner.Text()
		requestStr += line + "\r\n"
		if scanner.Text() == "" {
			break
		}
	}
	req, err := shiori.ParseRequest(requestStr)
	if err != nil {
		return dictionary.ResponseInternalServerError(), err
	}

	log.Println(req)

	// リクエストヘッダがなければ、統一的操作のために初期化しておきます。
	if req.Headers == nil {
		req.Headers = shiori.RequestHeaders{}
	}
	// ID ヘッダにはイベント名が入っています。
	// イベントに対応するハンドラが定義されていれば呼び出します。
	if event, ok := req.Headers["ID"]; ok {
		if handler, ok := Dictionary.Handlers[event]; ok {
			return handler(req, Dictionary.Variables)
		}
	}

	return dictionary.ResponseNoContent(), nil
}

func unload() error {
	if err := dictionary.SaveVariables(*Dictionary.Variables); err != nil {
		return err
	}

	if logFile != nil {
		logFile.Close()
	}

	return nil
}
