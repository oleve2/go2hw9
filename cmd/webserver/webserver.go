package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/wool/go2hw9/pkg/card"
)

//
func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

//
func execute() (err error) {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Println(cerr)
			if err == nil {
				err = cerr
			}
		}
	}()
	for {
		conn, err := listener.Accept() // для клиентов
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

//
func handle(conn net.Conn) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	r := bufio.NewReader(conn)
	const delim = '\n'
	line, err := r.ReadString(delim)
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		log.Printf("received: %s\n", line)
		return
	}
	log.Printf("received: %s\n", line)

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		log.Printf("invalid request line: %s", line)
		return
	}

	time.Sleep(time.Second * 2) //10
	path := parts[1]

	switch path {
	case "/":
		err = writeIndex(conn)
	case "/operations.csv":
		err = writeOperations(conn)
	case "/operations.json":
		err = writeOperationsJSON(conn)
	case "/operations.xml":
		err = writeOperationsXML(conn)
	default:
		err = write404(conn)
	}
	if err != nil {
		log.Println(err)
		return
	}
}

// Базовая страница
func writeIndex(writer io.Writer) error {
	username := "Василий"
	balance := "1 000.50" // FIXME: написать функцию-форматтер

	page, err := ioutil.ReadFile("web/template/index.html")
	if err != nil {
		return err
	}
	page = bytes.ReplaceAll(page, []byte("{username}"), []byte(username))
	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(balance))

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

// Выгрузка csv
func writeOperations(writer io.Writer) error {
	// TODO: Generate CSV
	c1 := card.InitCard()
	page, err := card.MakeCSV(c1.Transactions)
	if err != nil {
		log.Println(err)
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: text/csv",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

// writeOperationsJSON - экспорт json
func writeOperationsJSON(writer io.Writer) error {
	c1 := card.InitCard()
	v, err := card.MakeJSON(c1.Transactions)
	if err != nil {
		log.Println(err)
	}
	page := []byte(v)

	return writeResponse(writer, 200, []string{
		"Content-Type: application/json",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

// экспорт xml
func writeOperationsXML(writer io.Writer) error {
	c1 := card.InitCard()
	v, err := card.MakeXML(c1.Transactions)
	if err != nil {
		log.Println(err)
	}
	page := []byte(v)

	return writeResponse(writer, 200, []string{
		"Content-Type: application/xml",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

// Ответ - 404
func write404(writer io.Writer) error {
	page, err := ioutil.ReadFile("web/template/404.html")
	if err != nil {
		return err
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

//
func writeResponse(writer io.Writer, status int, headers []string, content []byte) error {
	const CRLF = "\r\n"
	var err error

	w := bufio.NewWriter(writer)
	_, err = w.WriteString(fmt.Sprintf("HTTP/1.1 %d OK%s", status, CRLF))
	if err != nil {
		return err
	}

	for _, h := range headers {
		_, err = w.WriteString(h + CRLF)
		if err != nil {
			return err
		}
	}

	_, err = w.WriteString(CRLF)
	if err != nil {
		return err
	}
	_, err = w.Write(content)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}
