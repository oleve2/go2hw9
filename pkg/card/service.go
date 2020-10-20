package card

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Card - тип единичной карты
type Card struct {
	ID           int64
	Type         string
	BankName     string
	CardNumber   string
	CardDueDate  string
	Balance      int64
	Transactions []*Transaction
}

// Transaction - Transaction
type Transaction struct {
	XMLName  string `xml:"transaction"`              //
	ID       int64  `json:"id" xml:"id"`             //
	TranType string `json:"trantype" xml:"trantype"` //
	TranSum  int64  `json:"transum" xml:"transum"`   //
	TranDate int64  `json:"trandate" xml:"trandate"` //  unix timestamp
	MccCode  string `json:"mcccode" xml:"mcccode"`   //
	Status   string `json:"status" xml:"status"`     //
	OwnerID  int64  `json:"ownerid" xml:"ownerid"`   //
}

// Transactions -
type Transactions struct {
	XMLName      string         `xml:"transactions"`
	Transactions []*Transaction `xml:"transaction"`
}

// Service - итоговый сервис для передачи во вне
type Service struct {
	BankName string
	Cards    []*Card
}

// SearchByNumber - поиск карты в сервисе по её номеру (возвращаем её номер в списке карт)
func (s *Service) SearchByNumber(number string) (*Card, bool) {
	for _, card := range s.Cards {
		if card.CardNumber == number {
			return card, true
		}
	}
	return nil, false
}

// AddTransaction - добавление транцакции
func AddTransaction(card *Card, transaction *Transaction) {
	card.Transactions = append(card.Transactions, transaction)
}

// функция проверки вхождения
func valInSlice(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// SumByMCC - функция сумм по коду mmc
func SumByMCC(transactions []*Transaction, mcc []string) int64 {
	var totalMcc int64
	for _, v := range transactions {
		if valInSlice(v.MccCode, mcc) == true {
			totalMcc += v.TranSum
		}
	}
	return totalMcc
}

// PrintCardTrans -
func PrintCardTrans(c *Card) {
	for _, v := range c.Transactions {
		fmt.Println(v.ID, v.TranDate, v.TranSum, v.TranType, v.MccCode, v.Status)
	}
}

// SortSlice -
func SortSlice(c *Card, asc bool) {
	if asc == true {
		sort.SliceStable(c.Transactions, func(i, j int) bool { return c.Transactions[i].TranSum < c.Transactions[j].TranSum })
	} else {
		sort.SliceStable(c.Transactions, func(i, j int) bool { return c.Transactions[i].TranSum > c.Transactions[j].TranSum })
	}
}

// Sum - подсчет суммы транзакций
func Sum(transactions []int64) int64 {
	var res int64 = 0
	for _, v := range transactions {
		res += v
	}
	return res
}

// MakeTransMap -
func MakeTransMap(trans []*Transaction) map[string][]int64 {
	var mp = make(map[string][]int64)
	for _, v := range trans {
		var y string = strconv.Itoa(time.Unix(v.TranDate, 0).UTC().Year())
		var m string = strconv.FormatInt(int64(time.Unix(v.TranDate, 0).UTC().Month()), 10)

		var key string
		if len(m) == 1 {
			key = y + " 0" + m
		} else if len(m) == 2 {
			key = y + " " + m
		}
		mp[key] = append(mp[key], v.TranSum)
	}
	return mp
}

// SumConcurrently -
func SumConcurrently(trans []*Transaction, goroutines int) int64 {
	transMap := MakeTransMap(trans)

	lenTM := len(transMap)
	wg := sync.WaitGroup{}
	wg.Add(lenTM)

	total := int64(0)
	var sumByMonths = make(map[string]int64)
	mx := sync.Mutex{}

	for i, v := range transMap {
		yyyymm := i
		trans := v
		go func() {
			sum := Sum(trans)
			mx.Lock()
			sumByMonths[yyyymm] = sum
			total += sum
			mx.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	for k, v := range sumByMonths {
		fmt.Printf("%v : %d\n", k, v)
	}
	fmt.Println(sumByMonths)
	return total
}

/*
F1) Обычная функция, которая принимает на вход слайс транзакций и id владельца
- возвращает map с категориями и тратами по ним (сортировать они ничего не должна)

F2) Функция с mutex'ом, который защищает любые операции с map, соответственно, её задача: разделить слайс транзакций
на несколько кусков и в отдельных горутинах посчитать map'ы по кускам, после чего собрать всё в один большой map.
Важно: эта функция внутри себя должна вызывать функцию из п.1

F3) Функция с каналами, соответственно, её задача: разделить слайс транзакций на несколько кусков и в отдельных
горутинах посчитать map'ы по кускам, после чего собрать всё в один большой map (передавайте рассчитанные куски по каналу).
Важно: эта функция внутри себя должна вызывать функцию из п.1

F4) Функция с mutex'ом, который защищает любые операции с map, соответственно, её задача: разделить слайс транзакций
на несколько кусков и в отдельных горутинах посчитать, но теперь горутины напрямую пишут в общий map с результатами.
Важно: эта функция внутри себя не должна вызывать функцию из п.1

*/

// DiviveTranSlcToParts - разделить транзакции на NumberOfParts частей
func DiviveTranSlcToParts(tr []*Transaction, NumberOfParts int64) map[int64][]*Transaction {
	mp := make(map[int64][]*Transaction)
	slcLen := int64(len(tr))
	var partSize int64
	if slcLen%NumberOfParts == 0 {
		partSize = slcLen / NumberOfParts
	} else {
		partSize = slcLen/NumberOfParts + 1
	}

	var start, finish int64
	start = 0
	for i := 0; i < int(NumberOfParts); i++ {
		finish = start + partSize
		if finish < slcLen {
			mp[int64(i)] = tr[int(start):int(finish)]
		} else {
			mp[int64(i)] = tr[int(start):]
			return mp
		}
		start = finish
	}
	return mp
}

// F1 - сумма в лоб
func F1(tr []*Transaction, ownerID int64) map[string]int64 {
	mp := make(map[string]int64)
	for _, v := range tr {
		//fmt.Printf("code %v, code name %v, transum %d \n", v.MccCode, TranslateMCC(v.MccCode), v.TranSum)
		if v.OwnerID == ownerID {
			mp[TranslateMCC(v.MccCode)] += v.TranSum
		}
	}
	//fmt.Println(mp)
	return mp
}

// F2 - сумма конкурентно через мьютексы
func F2(tr []*Transaction, ownerID int64) map[string]int64 {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	result := make(map[string]int64)

	transSplit := DiviveTranSlcToParts(tr, 100)

	for _, v := range transSplit { // TODO здесь ваши условия разделения
		wg.Add(1)
		part := v
		go func() {
			m := F1(part, ownerID) // Categorize(part)
			mu.Lock()
			// TODO: вы перекладываете данные из m в result
			// TODO: подсказка - сделайте цикл по одной из map и смотрите, есть ли такие ключи в другой, если есть - прибавляйте
			for k, v := range m {
				result[k] += v
			}
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return result
}

/*
Функция с каналами, соответственно, её задача: разделить слайс транзакций на несколько кусков и в отдельных горутинах посчитать map'ы по кускам,
после чего собрать всё в один большой map (передавайте рассчитанные куски по каналу).
Важно: эта функция внутри себя должна вызывать функцию из п.1
*/

// F3 - конкуретноый подсчет через каналы
func F3(tr []*Transaction, ownerID int64) map[string]int64 {
	result := make(map[string]int64)
	ch := make(chan map[string]int64)

	transSplit := DiviveTranSlcToParts(tr, 100)

	for _, v := range transSplit { // TODO здесь ваши условия разделения
		part := v // transactions[x:y]
		go func(ch chan<- map[string]int64) {
			ch <- F1(part, ownerID) //Categorize(part)
		}(ch)
	}

	partsCount := len(transSplit)
	finished := 0
	for value := range ch { // range result
		// TODO: вы перекладываете данные из m в result
		// TODO: подсказка - сделайте цикл по одной из map и смотрите, есть ли такие ключи в другой, если есть - прибавляйте
		for k, v := range value {
			result[k] += v
		}
		finished++
		if finished == partsCount {
			break
		}
	}
	return result
}

/*
F4) Функция с mutex'ом, который защищает любые операции с map, соответственно, её задача: разделить слайс транзакций
на несколько кусков и в отдельных горутинах посчитать, но теперь горутины напрямую пишут в общий map с результатами.
Важно: эта функция внутри себя не должна вызывать функцию из п.1
*/
// F4 -
func F4(tr []*Transaction, ownerID int64) map[string]int64 {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	result := make(map[string]int64)

	transSplit := DiviveTranSlcToParts(tr, 100)

	for _, v := range transSplit { // TODO здесь ваши условия разделения
		wg.Add(1)
		part := v //transactions[x:y]
		go func() {
			for _, t := range part {
				// TODO: 1. берём конкретную транзакцию
				// TODO: 2. смотрим, подходит ли по id владельца
				if t.OwnerID == ownerID {
					mu.Lock()
					// TODO: 3. если подходит, то закидываем в общий `map`
					result[TranslateMCC(t.MccCode)] += t.TranSum
					mu.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return result
}

/*
Экспорт и импорт транзакций (csv)

0) База
на вход - объект card с транзакцийми; путь экспорта + назв.файла;
на выход - ошибка (если её нет то nil)

1) Поля:
ID       int64
TranType string
TranSum  int64
TranDate int64 // unix timestamp
MccCode  string
Status   string
OwnerID  int64

2) План реализации:
) на вход - объект карт

*/

// MapRowToTransaction - конвертация в тип Transaction
func MapRowToTransaction(s [][]string) []*Transaction {
	trans := make([]*Transaction, 0)
	for _, v := range s {
		id2, _ := strconv.ParseInt(v[0], 10, 64)
		transum2, _ := strconv.ParseInt(v[2], 10, 64)

		layout := "2006-01-02 15:04:05 +0300 MSK"
		trandate2, _ := time.Parse(layout, v[3]) //"2014-11-12T11:45:26.371Z"
		trandate3 := trandate2.Unix()

		owner2, _ := strconv.ParseInt(v[6], 10, 64)

		tr := &Transaction{
			"",
			id2,       //ID       int64
			v[1],      //TranType string
			transum2,  //TranSum  int64
			trandate3, //TranDate int64 // unix timestamp
			v[4],      //MccCode  string
			v[5],      //Status   string
			owner2,    //OwnerID  int64
		}
		trans = append(trans, tr)
	}
	return trans
}

// ExportToCSV - export to csv
func ExportToCSV(tr []*Transaction, exportPath string) error {
	if len(tr) == 0 {
		return nil
	}

	records := make([][]string, 0)
	for _, v := range tr {
		record := []string{
			strconv.FormatInt(v.ID, 10),
			v.TranType,
			strconv.FormatInt(v.TranSum, 10),
			time.Unix(v.TranDate, 0).String(), // TranDate
			v.MccCode,
			v.Status,
			strconv.FormatInt(v.OwnerID, 10),
		}
		records = append(records, record)
	}

	file, err := os.Create(exportPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func(c io.Closer) {
		if err := c.Close(); err != nil {
			log.Println(err)
		}
	}(file)

	w := csv.NewWriter(file)
	w.WriteAll(records)

	return nil
}

// ImportFromCSV - import from csv
func ImportFromCSV(importPath string) ([]*Transaction, error) {
	file, err := os.Open(importPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println(cerr)
		}
	}(file)

	reader := csv.NewReader(file)
	records := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return nil, err
			}
			//records = append(records, record)   <-- нужно ли прикреплять EOF к концу слайса? наверно нет ...
			break
		}
		records = append(records, record)
	}

	transConv := MapRowToTransaction(records)

	return transConv, nil
}

// ExporttoJSON - export to json
func ExporttoJSON(tr []*Transaction, exportPath string) error {
	if len(tr) == 0 {
		return nil
	}
	encData, err := json.Marshal(tr)
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(exportPath, encData, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// ImportFromJSON - import from json
func ImportFromJSON(importPath string) ([]*Transaction, error) {
	file, err := os.Open(importPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println(cerr)
		}
	}(file)

	content, err := ioutil.ReadFile(importPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var decoded []*Transaction
	err = json.Unmarshal(content, &decoded) // важно: передаём указатель
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println(reflect.TypeOf(decoded), decoded)
	return decoded, nil
}

// ExportXML -
func ExportXML(tr []*Transaction, exportPath string) error {
	if len(tr) == 0 {
		return nil
	}

	trs := &Transactions{Transactions: tr}
	encData, err := xml.Marshal(trs)
	if err != nil {
		log.Println(err)
		return err
	}
	encData = append([]byte(xml.Header), encData...)
	err = ioutil.WriteFile(exportPath, encData, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// ImportXML -
func ImportXML(importPath string) ([]*Transaction, error) {
	file, err := os.Open(importPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println(cerr)
		}
	}(file)

	content, err := ioutil.ReadFile(importPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	//
	var decoded Transactions
	err = xml.Unmarshal(content, &decoded)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Printf("%#v", decoded)
	trans := decoded.Transactions

	return trans, nil
}

// MakeCSV - HW9 - make csv
/*
func MakeCSV(tr []*Transaction) []byte {
	if len(tr) == 0 {
		return nil
	}
	tot := make([]string, 0)
	//records := make([][]string, 0)
	for _, v := range tr {
		record := []string{
			strconv.FormatInt(v.ID, 10),
			v.TranType,
			strconv.FormatInt(v.TranSum, 10),
			time.Unix(v.TranDate, 0).String(), // TranDate
			v.MccCode,
			v.Status,
			strconv.FormatInt(v.OwnerID, 10),
		}
		tmp1 := strings.Join(record[:], ",")
		tmp1 = append(tmp1, "\n")
		//records = append(records, record)
	}

	return []byte(records)
}*/

// MakeJSON -
func MakeJSON(tr []*Transaction) ([]byte, error) {
	if len(tr) == 0 {
		return nil, nil
	}
	encData, err := json.Marshal(tr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return encData, nil
}

// MakeXML -
func MakeXML(tr []*Transaction) ([]byte, error) {
	if len(tr) == 0 {
		return nil, nil
	}
	trs := &Transactions{Transactions: tr}
	encData, err := xml.Marshal(trs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	encData = append([]byte(xml.Header), encData...)
	return encData, nil
}
