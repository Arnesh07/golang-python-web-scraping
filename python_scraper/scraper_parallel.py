import concurrent.futures
import csv
import requests
import time
from time import sleep
import threading
from bs4 import BeautifulSoup

thread_local = threading.local()
max_workers = 5
tickers_to_be_scraped = 500

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
                ticker["url"] = "https://finance.yahoo.com/quote/" + row["Symbol"]
                tickers.append(ticker)
            line_count += 1
    return tickers

def get_session():
    if not hasattr(thread_local, "session"):
        thread_local.session = requests.Session()
    return thread_local.session

# function to fetch single quote and parse data
def fetch_quote(ticker):
    session = get_session()
    headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36'}
    session.headers = headers
    
    with session.get(ticker["url"]) as response:
        print(ticker["url"])
        
        results = BeautifulSoup(response.content, "html.parser")
        quote_price = results.find("span", class_="Trsdu(0.3s) Fw(b) Fz(36px) Mb(-4px) D(ib)")
        
        if quote_price == None:
            print("Not fetched price")
            return
        
        ticker["quote_price"] = float(quote_price.text)
        print(quote_price.text)

# function that creates a thread pool executor
def fetch_quotes(tickers):
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        results = executor.map(fetch_quote, tickers)
        executor.shutdown(wait=True)


#######################MAIN######################

start_time = time.time()

# Read csv file to obtain ticker information
tickers = read_csv("../nasdaq_screener_1635280898552.csv")
tickers = tickers[:tickers_to_be_scraped]
fetch_quotes(tickers)
count = 0

# For loop to check the number of tickers for which the data has been
# successfully scraped
for ticker in tickers:
    if "quote_price" in ticker:
        count += 1
print("Data Successfully Scraped For:")
print(count, "/", tickers_to_be_scraped)

print("Code execution time: %s seconds" % (time.time() - start_time))