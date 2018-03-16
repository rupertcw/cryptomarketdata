package main

import "time"

const (
	CurrencySymbol              = "USD"
	DBConnectionString          = "user=%s dbname=%s sslmode=disable"
	DBDriverName                = "postgres"
	DBName                      = "crypto"
	DBTableSuffix               = "dayohlc"
	DBUser                      = "postgres"
	HistoricalDayAPIURL         = "https://min-api.cryptocompare.com/data/histoday"
	HistoricalDayAPIQueryString = "%s?fsym=%s&tsym=%s"
	HTTPClientTimeout           = 10 * time.Second
)