package config

import (
    "os"
    "log"
    "database/sql"
    "github.com/go-sql-driver/mysql"
)

var DB *sql.DB 

func InitializeDB() {
    dbConfig := mysql.Config{
        User: os.Getenv("DB2_USERNAME"),
        Passwd: os.Getenv("DB2_PASSWORD"),
        Net: "tcp",
        Addr: os.Getenv("DB2_HOST")+":"+os.Getenv("DB2_PORT"),
        DBName: os.Getenv("DB2_NAME"),
        AllowNativePasswords: true,
    }

    var err error
    DB, err = sql.Open("mysql", dbConfig.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := DB.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }

    log.Println("Database is Connected")
}
