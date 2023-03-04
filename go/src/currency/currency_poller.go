package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "database/sql"

    _ "github.com/denisenkom/go-mssqldb"
)

type Candle struct {
    Complete bool `json:"complete"`
    Volume   int  `json:"volume"`
    Time     string `json:"time"`
    Mid      struct {
        Open     float64 `json:"o,string"`
        High     float64 `json:"h,string"`
        Low      float64 `json:"l,string"`
        Close    float64 `json:"c,string"`        
    } `json:"mid"`
}

type ApiResponse struct {
    Instrument   string   `json:"instrument"`
    Granularity  string   `json:"granularity"`
    Candles      []Candle `json:"candles"`
}

const apiKey           = "dc0b7e333a8ea31e0739b92e835af2c1-31e369bf78ec5d9df867fedacc2097d5"
const baseURL          = "https://api-fxpractice.oanda.com/v3"
const instrumentsPath  = "/instruments/"
const query            = "?granularity=D&price=A&"
const accountID        = "101-001-5466925-001"

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

    for {
        inst := "USD_JPY"
        count := 12
        apiResponse := getData(inst, count)

        fmt.Println("Instrument:", apiResponse.Instrument)
        fmt.Println("Granularity:", apiResponse.Granularity)
        fmt.Println("Candles:")
        for _, candle := range apiResponse.Candles {
            // Insert the candlestick data into the database
            err = insertCandlestick(db, candle, apiResponse.Instrument)
            if err != nil {
                fmt.Println("Error inserting candlestick data:", err)
                continue
            }

            fmt.Printf("  Time: %s, Open: %s, High: %s, Low: %s, Close: %s\n", candle.Time, candle.Mid.Open, candle.Mid.High, candle.Mid.Low, candle.Mid.Close)
        }

        time.Sleep(30 * time.Second)
    }
}


func insertCandlestick(db *sql.DB, candle Candle, instrument string) error {
    // Convert the time string to a time.Time object
    t, err := time.Parse(time.RFC3339Nano, candle.Time)
    if err != nil {
        return fmt.Errorf("error parsing time: %s", err)
    }


    // Insert some data into the table using MERGE INTO
    stmt, err := db.Prepare(`
        MERGE INTO candlesticks AS target
        USING (VALUES (?, ?, ?, ?, ?, ?, ?, ?)) AS source (instrument, bid_time, bid_open, bid_high, bid_low, bid_close, volume, complete)
        ON target.instrument = source.instrument AND target.bid_time = source.bid_time
        WHEN NOT MATCHED BY target THEN
            INSERT (instrument, bid_time, bid_open, bid_high, bid_low, bid_close, volume, complete)
            VALUES (source.instrument, source.bid_time, source.bid_open, source.bid_high, source.bid_low, source.bid_close, source.volume, source.complete);
    `)
    if err != nil {
        return fmt.Errorf("error preparing merge into statement data: %s", err)
    }
    defer stmt.Close()

    // Execute the SQL statement to insert the candlestick data
    _, err = stmt.Exec(instrument, t, candle.Mid.Open, candle.Mid.High, candle.Mid.Low, candle.Mid.Close, candle.Volume, candle.Complete)
    if err != nil {
        return fmt.Errorf("error inserting data: %s", err)
    }

    // Print the candlestick data
    fmt.Printf("Time: %s, Open: %s, High: %s, Low: %s, Close: %s, Volume: %s, Complete: %t\n",
        candle.Time, candle.Mid.Open, candle.Mid.High, candle.Mid.Low, candle.Mid.Close, candle.Volume, candle.Complete)

    return nil
}

func getData(instrument string, count int) ApiResponse {

    completeURL   := baseURL + instrumentsPath + instrument + "/candles"
    granularity   := "S10"
    price         := "M"
    authorization := "Bearer "+apiKey

    req, _ := http.NewRequest("GET", completeURL, nil)
    req.Header.Set("Authorization", authorization)
    req.Header.Set("Content-Type", "application/json")

    q := req.URL.Query()
    q.Add("count", fmt.Sprintf("%v", count))
    q.Add("price", price)
    q.Add("granularity", granularity)
    req.URL.RawQuery = q.Encode()

    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    var apiResponse ApiResponse
    json.Unmarshal([]byte(body), &apiResponse)

    return apiResponse
}
