# Golang v/s Python for Web Scraping

This repository contains code of web scrapers in both languages: Golang and Python. Yahoo Finance is used to scrape ticker quotes.


## Golang

The go_scraper contains 5 files:
- scraper_seq.go:
This contains scraper that runs sequentially and uses Gocolly framework for scraping.
- scraper_par_gocolly_queue.go:
Scraper that uses gocolly framework and the gocolly's queue mechanism to run scraping on multiple worker threads.
- scraper_par_gocolly_goroutine.go:
Scraper that uses gocolly framework and the uses goroutines to manually manage threads.
- scraper_par_gocolly_parallelism.go:
Scraper that uses gocolly framework and the uses gocolly's parallel mechanism to manage parallelism for scraping.
- scraper_par_goquery.go:
This scraper uses the goquery package for web scraping, and manages worker threads using goroutines manually. This method produces the fastest results.


## Python

The python_scraper contains two files:
- scraper_seq.go:
Sequential code that uses BeautifulSoup for scraping.
- scraper_par.go:
Parallel code that uses BeautifulSoup for scraping and uses Python's ThreadPoolExecutor for managing worker threads.

## Results

