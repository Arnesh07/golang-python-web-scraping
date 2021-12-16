import csv
import requests
import time
from bs4 import BeautifulSoup

def read_csv(filepath):
    tickers = []
    with open(filepath, mode='r') as csv_file:
        csv_reader = csv.DictReader(csv_file)
        line_count = 0
        for row in csv_reader:
            if line_count == 0:
                line_count += 1
                continue
            else:
                ticker = {}
                ticker["name"] = row["Name"]
                ticker["symbol"] = row["Symbol"]
                tickers.append(ticker)
            line_count += 1
    return tickers

#######################MAIN######################

start_time = time.time()

# Read csv file to obtain ticker information
tickers = read_csv("../nasdaq_screener_1635280898552.csv")

count = 0

# For loop to fetch every ticker data.
for ticker in tickers:
    count += 1
    URL = "https://finance.yahoo.com/quote/" + str(ticker["symbol"])
    print(URL)
    headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36'}
    
    page = requests.get(URL, headers=headers)
    results = BeautifulSoup(page.content, "html.parser")
    quote_price = results.find("span", class_="Trsdu(0.3s) Fw(b) Fz(36px) Mb(-4px) D(ib)")
    
    if quote_price == None:
        continue
    ticker["quote_price"] = quote_price.text
    print(quote_price.text)
    if count == 100:
        break

count = 0
for ticker in tickers:
    if "quote_price" in ticker:
        count += 1
print("Data Successfully Scraped For:")
print(count, "/", 100)


print("Code execution time: %s seconds" % (time.time() - start_time))