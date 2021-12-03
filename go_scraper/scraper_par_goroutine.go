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
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

type Ticker struct {
	Name   string
	Symbol string
	Price  float64
	URL    string
}

//slice to store symbols of tickers
var tickerSymbols = make([]string, 0)

//max number of goroutines to be spawned at a time
var maxGoroutines int = 10
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
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return tickers
}

func main() {
	start := time.Now()
	// reads csv file to get stock tickers.
	tickers := readCsvFile("../nasdaq_screener_1635280898552.csv")

	c := colly.NewCollector()

	// waitgroup is used to manage the goroutines
	var wg sync.WaitGroup

	//guard is used to ensure that the number of goroutines at a time is limited
	guard := make(chan struct{}, maxGoroutines)

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

		tickers[symbol] = ticker
		fmt.Println(price)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		<-guard
		wg.Done()
	})

	c.OnResponse(func(r *colly.Response) {
		<-guard
		wg.Done()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	sort.Strings(tickerSymbols)
	for i := 0; i < tickersToBeScraped; i++ {
		guard <- struct{}{} // would block if guard channel is already filled
		wg.Add(1)
		go c.Visit(tickers[tickerSymbols[i]].URL)
	}

	wg.Wait()
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
