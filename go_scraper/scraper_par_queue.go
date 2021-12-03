package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

type Ticker struct {
	Name   string
	Symbol string
	Price  float64
	URL    string
}

//slice to store symbols of tickers
var tickerSymbols = make([]string, 0)

// Number of consumer threads to be spawned simultaneously
var consumerThreads int = 10

var tickersToBeScraped int = 2000

func readCsvFile(filePath string) map[string]Ticker {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	tickers := make(map[string]Ticker)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		url := "https://finance.yahoo.com/quote/" + record[0]
		var ticker = Ticker{
			Name:   record[1],
			Symbol: record[0],
			URL:    url,
		}
		tickerSymbols = append(tickerSymbols, record[0])
		tickers[record[0]] = ticker
	}
	// records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return tickers
}

func main() {
	start := time.Now()
	// reads csv file to get stock tickers.

	// create a request queue with 5 consumer threads
	tickerQueue, _ := queue.New(
		consumerThreads, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	tickers := readCsvFile("../nasdaq_screener_1635280898552.csv")

	c := colly.NewCollector()

	c.OnHTML("#quote-header-info", func(e *colly.HTMLElement) {
		// parsing data and storing it in a map
		name := e.ChildText("h1")
		quote := e.ChildTexts("span")

		temp := strings.Split(name, "(")
		name = temp[0]
		symbol := temp[1][:len(temp[1])-1]

		price, _ := strconv.ParseFloat(quote[3], 32)
		price = math.Round(price/0.01) * 0.01

		ticker := tickers[symbol]
		ticker.Name = name
		ticker.Price = price

		fmt.Println(price)
		tickers[symbol] = ticker
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	sort.Strings(tickerSymbols)
	for i := 0; i < tickersToBeScraped; i++ {
		// Add URLs in a queue
		tickerQueue.AddURL(tickers[tickerSymbols[i]].URL)
	}
	tickerQueue.Run(c)

	// For loop which checks the number of tickers for which the data has been
	// successfully scraped
	count := 0
	for i := 0; i < tickersToBeScraped; i++ {
		if tickers[tickerSymbols[i]].Price != 0.00 {
			count++
		}
	}
	fmt.Println("Data Successfully Scraped For:")
	fmt.Println(count, "/", tickersToBeScraped)
	fmt.Println(time.Since(start))

}
