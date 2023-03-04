package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/denisenkom/go-mssqldb"
)

type Candlestick struct {
    Instrument string
    Time       time.Time
    Open       float64
    High       float64
    Low        float64
    Close      float64
    Volume     int
    Complete   bool
}


var db *sql.DB
var server = "currency-db.database.windows.net"
var port = 1433
var user = "jeffolen4"
var password = "zua.GJF6eqy1wqc9fjp"
var database = "currency_db"

func main() {

    // Connect to the Azure database
    connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
        server, user, password, port, database)

    db, err := sql.Open("mssql",connString)
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return
    }
    defer db.Close()

    // Query the candlesticks table.
    queryString := fmt.Sprintf("SET ROWCOUNT 18 select instrument, bid_time, bid_open, bid_high, bid_low, bid_close, volume, complete FROM candlesticks where instrument = '%s' ORDER BY bid_time DESC", "USD_JPY")
    rows, err := db.Query(queryString)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Iterate over the rows and create Candlestick structs.
    var candlesticks []Candlestick
    for rows.Next() {
        var candlestick Candlestick
        err := rows.Scan(&candlestick.Instrument, &candlestick.Time, &candlestick.Open, &candlestick.High, &candlestick.Low, &candlestick.Close, &candlestick.Volume, &candlestick.Complete)
        if err != nil {
            log.Fatal(err)
        }
        candlesticks = append(candlesticks, candlestick)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    // Analyze the candlesticks.
    priceIncreasing, priceDecreasing, volumeIncreasing, volumeDecreasing := analyzeCandlesticks(candlesticks)

    // Print out the results.
    fmt.Printf("Price increasing: %v\n", priceIncreasing)
    fmt.Printf("Price decreasing: %v\n", priceDecreasing)
    fmt.Printf("Volume increasing: %v\n", volumeIncreasing)
    fmt.Printf("Volume decreasing: %v\n", volumeDecreasing)
    fmt.Printf("instrument: %s\n", candlesticks[0].Instrument)
}

func analyzeCandlesticks(candlesticks []Candlestick) (bool, bool, bool, bool) {
    // Determine if the price and volume are increasing or decreasing.
    var priceIncreasing, priceDecreasing, volumeIncreasing, volumeDecreasing bool
    var totalPriceChange float64

    var avg_price_newest float64
    var avg_price_oldest float64
    var avg_volume_newest int
    var avg_volume_oldest int

    midway := int(len(candlesticks)/2)-1

    totalPriceChange = candlesticks[0].Close - candlesticks[len(candlesticks)-1].Close  


    // get early (first half) average
    for i := 0; i < midway+1; i++ {
        avg_volume_newest += candlesticks[i].Volume
        avg_price_newest  += candlesticks[i].Close
        fmt.Printf("%v - close: %v, volume: %v\n", i, candlesticks[i].Close, candlesticks[i].Volume)
    }
    avg_price_newest  = avg_price_newest/float64(midway+1)
    avg_volume_newest = avg_volume_newest/(midway+1)

    // get early (first half) average
    for i := int(len(candlesticks)/2); i < len(candlesticks); i++ {
        avg_volume_oldest += candlesticks[i].Volume
        avg_price_oldest  += candlesticks[i].Close
        fmt.Printf("%v - close: %v, volume: %v\n", i, candlesticks[i].Close, candlesticks[i].Volume)
    }
    avg_price_oldest  = avg_price_oldest/float64(midway+1)
    avg_volume_oldest = avg_volume_oldest/(midway+1)


    threshold := .0020 // 20 pips

    if candlesticks[0].Instrument == "USD_JPY" {
        threshold = threshold * 100
    }

    fmt.Printf("totalPriceChange: %v\n", totalPriceChange)
    fmt.Printf("threshold: %v\n", threshold)

    if totalPriceChange >= threshold && avg_price_newest > avg_price_oldest { // 20 pips
        priceIncreasing = true
    } else if totalPriceChange <= (threshold * -1) && avg_price_newest < avg_price_oldest { // -20 pips
        priceDecreasing = true
    }

    if avg_volume_newest > avg_volume_oldest {
        volumeIncreasing = true
    } else if avg_volume_newest < avg_volume_oldest {
        volumeDecreasing = true
    }

    return priceIncreasing, priceDecreasing, volumeIncreasing, volumeDecreasing
}