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

func initCard() *card.Card {
	card1 := &card.Card{ID: 1, Type: "Master", BankName: "Citi", CardNumber: "1111 2222 3333 4444", Balance: 20_000_00, CardDueDate: "2030-01-01",
		Transactions: []*card.Transaction{
			&card.Transaction{ID: 1, TranType: "purchase", OwnerID: 2, TranSum: 1735_55, TranDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5411", Status: "done Супермаркеты"},
			&card.Transaction{ID: 2, TranType: "purchase", OwnerID: 2, TranSum: 2000_00, TranDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5411", Status: "done"},
			&card.Transaction{ID: 3, TranType: "purchase", OwnerID: 2, TranSum: 1203_91, TranDate: time.Date(2020, 2, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5411", Status: "done Рестораны"},
			&card.Transaction{ID: 4, TranType: "purchase", OwnerID: 2, TranSum: 3562_21, TranDate: time.Date(2020, 2, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "1111", Status: ""},
			&card.Transaction{ID: 5, TranType: "purchase", OwnerID: 2, TranSum: 1111_11, TranDate: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "1111", Status: ""},
			&card.Transaction{ID: 6, TranType: "purchase", OwnerID: 2, TranSum: 2222_22, TranDate: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "1111", Status: ""},
			&card.Transaction{ID: 7, TranType: "purchase", OwnerID: 2, TranSum: 6666_66, TranDate: time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "3333", Status: ""},
			&card.Transaction{ID: 8, TranType: "purchase", OwnerID: 2, TranSum: 4444_44, TranDate: time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "3333", Status: ""},
			&card.Transaction{ID: 9, TranType: "purchase", OwnerID: 2, TranSum: 5555_55, TranDate: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5555", Status: ""},
			&card.Transaction{ID: 10, TranType: "purchase", OwnerID: 2, TranSum: 3333_33, TranDate: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5411", Status: ""},
			&card.Transaction{ID: 11, TranType: "purchase", OwnerID: 2, TranSum: 3333_33, TranDate: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5555", Status: ""},
			&card.Transaction{ID: 12, TranType: "purchase", OwnerID: 2, TranSum: 3333_33, TranDate: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5555", Status: ""},
			&card.Transaction{ID: 13, TranType: "purchase", OwnerID: 2, TranSum: 3333_33, TranDate: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix(), MccCode: "5411", Status: ""},
		},
	}
	return card1
}

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
		handle(conn)
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
	page := []byte("xxxx,0001,0002,1592373247\nxxxx,0001,0002,1592373247\nxxxx,0001,0002,1592373247\nxxxx,0001,0002,1592373247\n")
	//c1 := initCard()

	return writeResponse(writer, 200, []string{
		"Content-Type: text/csv",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

// writeOperationsJSON - экспорт json
func writeOperationsJSON(writer io.Writer) error {
	c1 := initCard()
	v, _ := card.MakeJSON(c1.Transactions)
	page := []byte(v)

	return writeResponse(writer, 200, []string{
		"Content-Type: application/json",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

// экспорт xml
func writeOperationsXML(writer io.Writer) error {
	c1 := initCard()
	v, _ := card.MakeXML(c1.Transactions)
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
func writeResponse(
	writer io.Writer,
	status int,
	headers []string,
	content []byte,
) error {
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
