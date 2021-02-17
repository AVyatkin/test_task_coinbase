package main

import (
    "log"
    "net/url"
    "sync"
)

var instruments = []string{"ETH-BTC", "BTC-USD", "BTC-EUR"}

func main() {
    db := initDb()
    defer db.Close()

    if !checkTableExist(db) {
        createTable(db)
    }

    wg := sync.WaitGroup{}

    u := url.URL{Scheme: "wss", Host: *addr, Path: ""}
    log.Printf("connecting to %s", u.String())

    for _, instrument := range instruments {
        wg.Add(1)
        go worker(db, wg, u, instrument)
    }

    wg.Wait()
}
