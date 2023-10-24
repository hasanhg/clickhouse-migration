package config

import (
	"os"
	"strconv"
)

func getDatabasePort() int {
	const (
		DefaultPort = 25060
	)

	port := DefaultPort
	strPort := os.Getenv("POSTGRES_PORT")
	if strPort != "" {
		port, _ = strconv.Atoi(strPort)
		if port == 0 {
			port = DefaultPort
		}
	}

	return port
}

func getClickHousePort() int {
	const (
		DefaultPort = 9000
	)

	port := DefaultPort
	strPort := os.Getenv("CLICKHOUSE_PORT")
	if strPort != "" {
		port, _ = strconv.Atoi(strPort)
		if port == 0 {
			port = DefaultPort
		}
	}

	return port
}

func getMaxIdleConn() int {
	const (
		DefaultIdleConn = 32
	)

	idleConn := DefaultIdleConn
	strIdleConn := os.Getenv("POSTGRES_MAX_IDLE_CONN")
	if strIdleConn != "" {
		idleConn, _ = strconv.Atoi(strIdleConn)
		if idleConn == 0 {
			idleConn = DefaultIdleConn
		}
	}

	return idleConn
}

func getMaxOpenConn() int {
	const (
		DefaultOpenConn = 32
	)

	openConn := DefaultOpenConn
	strOpenConn := os.Getenv("POSTGRES_MAX_OPEN_CONN")
	if strOpenConn != "" {
		openConn, _ = strconv.Atoi(strOpenConn)
		if openConn == 0 {
			openConn = DefaultOpenConn
		}
	}

	return openConn
}

func getMaxLifetime() int {
	const (
		DefaultMaxLifetime = 30
	)

	maxLifetime := DefaultMaxLifetime
	strMaxLifetime := os.Getenv("POSTGRES_MAX_LIFETIME")
	if strMaxLifetime != "" {
		maxLifetime, _ = strconv.Atoi(strMaxLifetime)
		if maxLifetime == 0 {
			maxLifetime = DefaultMaxLifetime
		}
	}

	return maxLifetime
}

var Default = &config{
	ClickHouse: Database{
		IP:       os.Getenv("CLICKHOUSE_IP"),
		Port:     getClickHousePort(),
		User:     os.Getenv("CLICKHOUSE_USERNAME"),
		Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		Name:     os.Getenv("CLICKHOUSE_NAME"),
	},
	Database: Database{
		Driver:          "postgres",
		IP:              os.Getenv("POSTGRES_HOST"),
		Port:            getDatabasePort(),
		User:            os.Getenv("POSTGRES_USER"),
		Password:        os.Getenv("POSTGRES_PASSWORD"),
		Name:            os.Getenv("POSTGRES_DB"),
		ConnMaxIdle:     getMaxIdleConn(),
		ConnMaxOpen:     getMaxOpenConn(),
		ConnMaxLifetime: getMaxLifetime(),
		Debug:           false,
	},
}
