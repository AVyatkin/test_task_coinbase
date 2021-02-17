package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "strconv"
    "time"
)

type TickData struct {
    Symbol string `json:"product_id"`
    Bid    string `json:"best_bid"`
    Ask    string `json:"best_ask"`
}

type Tick struct {
    Timestamp int64   `json:"timestamp"`
    Symbol    string  `json:"symbol"`
    Bid       float64 `json:"bid"`
    Ask       float64 `json:"ask"`
}

var mysqlSource = "mysql_user:mysql_password@tcp(localhost:5000)/mysql_db"

func initDb() *sql.DB {
    db, err := sql.Open("mysql", mysqlSource)
    if err != nil {
        panic(err)
    }

    db.SetConnMaxLifetime(time.Minute * 3)
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10)

    return db
}

func checkTableExist(db *sql.DB) bool {
    r, err := db.Query("SHOW TABLES LIKE 'ticks'")
    if err != nil {
        panic(err.Error())
    }

    for r != nil && r.Next() {
        var f string
        err = r.Scan(&f)
        if err != nil {
            return false
        }
        if f == "ticks" {
            return true
        }
    }

    return false
}

func createTable(db *sql.DB) {
    _, err := db.Query(`create table ticks
(
    timestamp bigint       not null,
    symbol    varchar(255) not null,
    bid       double       not null,
    ask       double       not null
);`)
    if err != nil {
        panic(err.Error())
    }
}

func writeData(db *sql.DB, tick Tick) {
    insertString := `insert into ticks (timestamp, symbol, bid, ask) values (` +
        strconv.FormatInt(tick.Timestamp, 10) + `, '` +
        tick.Symbol + `', ` +
        fmt.Sprintf("%f", tick.Bid) + `, ` +
        fmt.Sprintf("%f", tick.Ask) + `)`
    insert, err := db.Query(insertString)
    if err != nil {
        panic(err.Error())
    }
    insert.Close()
}

func readData(db *sql.DB) {
    rows, err := db.Query("select timestamp, symbol, bid, ask from ticks")
    if err != nil {
        log.Printf("err: %#v\n\n", err.Error())
    }

    for rows != nil && rows.Next() {
        var tick Tick
        err = rows.Scan(&tick.Timestamp, &tick.Symbol, &tick.Bid, &tick.Ask)
        if err != nil {
            log.Println("error:" + err.Error())
        }
        log.Printf("tick: %#v\n\n", tick)
    }
}
