package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	defer measureTime(time.Now())

	filepath := "./logs2.txt"
	storage := NewStorageInMemory()

	// read logs from the file and save each balance in storage
	ReadTransactionLog(filepath, storage)

	// print the output
	storage.Report()
}

func ReadTransactionLog(filePath string, storage *StorageInMemory) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tran := ParseTransaction(line)

		storage.Save(tran)
	}

	file.Close()
}

type Transaction struct {
	Id        string
	Operation string
	Amount    float64
}

func ParseTransaction(log string) *Transaction {
	items := strings.Split(log, ",")

	var tran *Transaction
	{
		id := items[0]
		ope := items[1]

		amount, err := strconv.ParseFloat(items[2], 64)
		if err != nil {
			panic(err)
		}

		tran = &Transaction{
			Id:        id,
			Operation: ope,
			Amount:    amount,
		}
	}

	return tran
}

type StorageInMemory struct {
	storage map[string]float64
}

func (s StorageInMemory) Save(tran *Transaction) {
	s.ensureBalanceEntryExists(tran.Id)

	switch tran.Operation {
	case "deposit":
		s.storage[tran.Id] += tran.Amount
	case "withdraw":
		s.storage[tran.Id] -= tran.Amount
	}
}

func (s StorageInMemory) ensureBalanceEntryExists(id string) {
	if _, ok := s.storage[id]; !ok {
		s.storage[id] = 0.0
	}
}

func (s StorageInMemory) Report() {
	for id, balance := range s.storage {
		fmt.Printf("%s %.2f\n", id, balance)
	}
}

func NewStorageInMemory() *StorageInMemory {
	return &StorageInMemory{
		storage: map[string]float64{},
	}
}

func measureTime(start time.Time) {
	log.Printf("Time taken: %s\n", time.Since(start))
}
