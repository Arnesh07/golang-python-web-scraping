package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Ticker struct {
	Name   string
	Symbol string
	Price  float64
	URL    string
}

// map to store tickers
var tickers = make(map[string]Ticker)

//slice to store symbols of tickers
var tickerSymbols = make([]string, 0)

// Number of consumer threads to be spawned simultaneously
var maxGoroutines int = 5

var tickersToBeScraped int = 500

// waitgroup is used to manage the goroutines
var wg sync.WaitGroup

//guard is used to ensure that the number of goroutines at a time is limited
var guard chan struct{} = make(chan struct{}, maxGoroutines)

func readCsvFile(filePath string) map[string]Ticker {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
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

func fetchQuotePrice(URL string) {
	// Request the HTML page.
	fmt.Println("Visiting", URL)
	client := &http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("error is not nil", err)
		<-guard
		wg.Done()
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println(res.StatusCode, res.Status)
		<-guard
		wg.Done()
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		<-guard
		wg.Done()
		return
	}

	// parse the required div HTML element
	doc.Find("#quote-header-info").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		name := s.Find("h1").Text()
		temp := strings.Split(name, "(")
		name = temp[0]
		symbol := temp[1][:len(temp[1])-1]

		quote := s.Find("span").Eq(3).Text()
		price, _ := strconv.ParseFloat(quote, 32)
		price = math.Round(price/0.01) * 0.01

		ticker := tickers[symbol]
		ticker.Name = name
		ticker.Price = price

		tickers[symbol] = ticker
		fmt.Println(price)
	})
	<-guard
	wg.Done()
}

func main() {
	start := time.Now()
	// reads csv file to get stock tickers.
	tickers = readCsvFile("../nasdaq_screener_1635280898552.csv")

	sort.Strings(tickerSymbols)
	for i := 0; i < tickersToBeScraped; i++ {
		guard <- struct{}{} // would block if guard channel is already filled
		wg.Add(1)
		go fetchQuotePrice(tickers[tickerSymbols[i]].URL)
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
