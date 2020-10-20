package main

import (
	"log"
	"time"

	"github.com/wool/go2hw9/pkg/card"
)

func main() {
	// инициалзация карты и транзакций
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
	log.Println(card1)
	var csvFname, jsonFname, xmlFname string = "./data/csvfile.csv", "./data/jsonfile.json", "./data/xmlfile.xml"
	log.Println(csvFname, jsonFname, xmlFname)

	// блок CSV
	if true {
		if true {
			// экспорт csv
			err := card.ExportToCSV(card1.Transactions, csvFname)
			if err != nil {
				log.Println(err)
			}
			log.Println("export csv - done")
		}
		if true {
			// импорт csv
			card2 := &card.Card{ID: 1, Type: "UberCard", BankName: "VelBank", CardNumber: "2222 2222 2222 2222", Balance: 1_000_0_000_00, CardDueDate: "2100-01-01"}
			transConv, err := card.ImportFromCSV(csvFname)
			if err != nil {
				log.Println(err)
			}
			log.Println("transConv ", transConv)
			card2.Transactions = transConv
			card.PrintCardTrans(card2)
		}
	}

	// блок JSON
	if true {
		if true {
			// экспорт JSON
			err := card.ExporttoJSON(card1.Transactions, jsonFname)
			if err != nil {
				log.Println(err)
			}
			log.Println("json exported!")
		}

		if true {
			// импорт JSON
			card2 := &card.Card{ID: 1, Type: "UberCard", BankName: "VelBank", CardNumber: "2222 2222 2222 2222", Balance: 1_000_0_000_00, CardDueDate: "2100-01-01"}
			tran2, err := card.ImportFromJSON(jsonFname)
			if err != nil {
				log.Println(err)
			}
			log.Println(tran2, err)
			card2.Transactions = tran2
			card.PrintCardTrans(card2)
		}
	}

	// блок XML
	if true {
		if true {
			// экспорт XML
			err := card.ExportXML(card1.Transactions, xmlFname)
			if err != nil {
				log.Println(err)
			}
			log.Println("xml exported!")
		}

		if true {
			// импорт XML
			card2 := &card.Card{ID: 1, Type: "UberCard", BankName: "VelBank", CardNumber: "2222 2222 2222 2222", Balance: 1_000_0_000_00, CardDueDate: "2100-01-01"}
			tran2, err := card.ImportXML(xmlFname)
			if err != nil {
				log.Println(err)
			}
			log.Println(tran2, err)
			card2.Transactions = tran2
			card.PrintCardTrans(card2)
		}
	}
}
