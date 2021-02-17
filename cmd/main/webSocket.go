package main

import (
    "database/sql"
    "encoding/json"
    "flag"
    "github.com/gorilla/websocket"
    "log"
    "net/url"
    "strconv"
    "sync"
    "time"
)

var addr = flag.String("addr", "ws-feed.pro.coinbase.com", "http service address")

func worker(db *sql.DB, wg sync.WaitGroup, u url.URL, instrument string) {
    defer wg.Done()

    if instrument == "" {
        log.Println("Worker stopped: instrument must be not empty!")
        return
    }

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Println("Error: " + err.Error())
        return
    }
    defer c.Close()

    message := `{
                "type": "subscribe",
                "channels": [{ "name": "ticker", "product_ids": ["` + instrument + `"] }]
            }`
    err = c.WriteMessage(websocket.TextMessage, []byte(message))
    if err != nil {
        log.Println("write:", err)
        return
    }

    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read error:", err)
            return
        }
        log.Printf("recv: %s\n\n", message)

        tickData := TickData{}
        err = json.Unmarshal(message, &tickData)
        if err != nil {
            log.Println("unmarshal error: " + err.Error())
            continue
        }

        if tickData.Symbol == "" {
            continue
        }

        bid, _ := strconv.ParseFloat(tickData.Bid, 64)
        ask, _ := strconv.ParseFloat(tickData.Ask, 64)

        writeData(db, Tick{
            time.Now().UnixNano() / int64(time.Millisecond),
            tickData.Symbol,
            bid,
            ask,
        })
    }
}
