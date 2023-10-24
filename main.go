package main

import (
	"clickhouse-migrations/config"
	"clickhouse-migrations/database"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func initDatabase() {
	for {
		err := database.InitDB()
		if err == nil {
			log.Printf("Connected to database...")
			break
		}
		log.Printf("Database connection error: %+v, waiting 5 sec...", err)
		time.Sleep(time.Duration(5) * time.Second)
	}

	for {
		err := database.InitClickHouse()
		if err == nil {
			log.Printf("Connected to ClickHouse...")
			return
		}
		log.Printf("ClickHouse connection error: %+v, waiting 5 sec...", err)
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func init() {
	godotenv.Load()
}

func main() {
	// Just a workaround for glog flag parsing: https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})

	cfg, err := config.Load()
	if err != nil || cfg == nil {
		fmt.Printf("[FATAL] %s", err)
		return
	}

	initDatabase()

	database.MigrateJobs()
}
