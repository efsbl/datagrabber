package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const BASE_URL = "https://jsonplaceholder.typicode.com/users"

func main() {
	workers := flag.Int("w", 1, "number of goroutines")

	flag.Parse()

	tasks := make(chan uint64)

	go ReadData(tasks)

	var wg sync.WaitGroup
	wg.Add(*workers)
	results := make(chan TaskResponse)
	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < *workers; i++ {
		go func() {
			for id := range tasks {
				GetData(id, results)
			}
			wg.Done()
		}()
	}

	WriteData(results)

}

func GetData(id uint64, results chan TaskResponse) {
	response, err := http.Get(fmt.Sprintf("%s/%d", BASE_URL, id))
	if err != nil {
		log.Fatalln("error in get", err)
	}
	defer response.Body.Close()

	var user1 User
	if err := json.NewDecoder(response.Body).Decode(&user1); err != nil {
		log.Fatalln("error decoding response", err)
	}

	// simulate another api call
	response2, err := http.Get(fmt.Sprintf("%s/%d", BASE_URL, id))
	if err != nil {
		log.Fatalln("error in get", err)
	}
	defer response.Body.Close()

	var user2 User
	if err := json.NewDecoder(response2.Body).Decode(&user2); err != nil {
		log.Fatalln("error decoding response", err)
	}

	results <- TaskResponse{resp1: user1, resp2: user2}
}

type TaskResponse struct {
	resp1, resp2 User
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ReadData(tasks chan uint64) {
	file, err := os.Open("ids.csv")

	if err != nil {
		log.Fatalln("error opening file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalln("error reading file %w", err)
		}

		id, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			log.Fatalln("error parsing id %w", err)
		}

		tasks <- id
	}

	close(tasks)
}

func WriteData(results chan TaskResponse) {
	file, err := os.Create("users.csv")
	if err != nil {
		log.Fatalln("error creating file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Name", "Email"}
	if err := w.Write(headers); err != nil {
		log.Fatalln("error writing to file", err)
	}

	for r := range results {
		var values1, values2 []string
		values1 = append(values1, strconv.FormatUint(uint64(r.resp1.ID), 10), r.resp1.Name, r.resp1.Email)
		if err := w.Write(values1); err != nil {
			log.Fatalln("error writing to file", err)
		}

		values2 = append(values2, strconv.FormatUint(uint64(r.resp2.ID), 10), r.resp2.Name, r.resp2.Email)
		if err := w.Write(values2); err != nil {
			log.Fatalln("error writing to file", err)
		}

		if err := w.Write([]string{""}); err != nil {
			log.Fatalln("error writing to file", err)
		}

	}
}
