package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type OHLC struct {
	DateTime 	float64 `json:"time"`
	Open 		float64 `json:"open"`
	High 		float64 `json:"high"`
	Low 		float64 `json:"low"`
	Close 		float64 `json:"close"`
}

type HistoricalDailyOHLC struct {
	Data []OHLC
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func NewDBClient() sql.DB {
	log.Println(fmt.Sprintf("Connecting user:%s to database:%s ...", DBUser, DBName))

	dbInfo := fmt.Sprintf(DBConnectionString, DBUser, DBName)
	db, err := sql.Open(DBDriverName, dbInfo)
	checkErr(err)

	return *db
}

func jsonToOHLCs(response http.Response) []OHLC {
	body, err := ioutil.ReadAll(response.Body)
	checkErr(err)

	var historicalDailyOHLC HistoricalDailyOHLC
	err = json.Unmarshal(body, &historicalDailyOHLC)
	checkErr(err)

	return historicalDailyOHLC.Data
}

func getHistoricalDailyOHLC(httpClient http.Client, coinSymbol string) []OHLC {
	log.Println(fmt.Sprintf("Fetching %s market data from %s...", coinSymbol, HistoricalDayAPIURL))

	url := fmt.Sprintf(HistoricalDayAPIQueryString, HistoricalDayAPIURL, coinSymbol, CurrencySymbol)
	response, err := httpClient.Get(url)
	checkErr(err)

	defer response.Body.Close()

	return jsonToOHLCs(*response)
}

func bulkInsert(dbClient sql.DB, ohlcs []OHLC, coinSymbol string) {
	dbTableName := fmt.Sprintf("%s_%s", strings.ToLower(coinSymbol), DBTableSuffix)
	log.Println(fmt.Sprintf("Bulk inserting %s market data into %s...", coinSymbol, dbTableName))

	txn, err := dbClient.Begin()
	checkErr(err)
	defer txn.Commit()

	stmt, err := txn.Prepare(pq.CopyIn(dbTableName, "time", "open", "high", "low", "close"))
	checkErr(err)
	defer stmt.Close()
	defer stmt.Exec()

	for _, ohlc := range ohlcs {
		_, err = stmt.Exec(
			time.Unix(int64(ohlc.DateTime), 0),
			ohlc.Open,
			ohlc.High,
			ohlc.Low,
			ohlc.Close,
		)
		checkErr(err)
	}
}

func GetMarketDataForCoin(coinSymbol string, httpClient http.Client, dbClient sql.DB, workGroup *sync.WaitGroup){
	defer workGroup.Done()

	ohlcs := getHistoricalDailyOHLC(httpClient, coinSymbol)

	bulkInsert(dbClient, ohlcs, coinSymbol)
}

func main() {
	var waitGroup sync.WaitGroup

	httpClient := http.Client{Timeout: HTTPClientTimeout}
	dbClient := NewDBClient()
	defer dbClient.Close()

	for _, coinSymbol := range os.Args[1:] {
		waitGroup.Add(1)
		go GetMarketDataForCoin(coinSymbol, httpClient, dbClient, &waitGroup)
	}
	waitGroup.Wait()

	log.Println(fmt.Sprintf("Disconnecting user:%s from database:%s...", DBUser, DBName))
}

