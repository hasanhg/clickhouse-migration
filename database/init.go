package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	gorp "gopkg.in/gorp.v1"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"

	"clickhouse-migrations/config"
)

var (
	ch *gorm.DB
	db *gorp.DbMap
)

func InitDB() error {

	var (
		cfg         = config.Config
		ip          = cfg.Database.IP
		port        = cfg.Database.Port
		user        = cfg.Database.User
		pass        = cfg.Database.Password
		name        = cfg.Database.Name
		maxIdle     = cfg.Database.ConnMaxIdle
		maxOpen     = cfg.Database.ConnMaxOpen
		maxLifetime = cfg.Database.ConnMaxLifetime
		debug       = cfg.Database.Debug
		err         error
		conn        *sql.DB
	)

	conn, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", ip, port, user, pass, name))

	if err != nil {
		return fmt.Errorf("Database connection error: %s", err.Error())
	}

	conn.SetMaxIdleConns(maxIdle)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	if err = conn.Ping(); err != nil {
		return fmt.Errorf("Database connection error: %s", err.Error())
	}

	db = &gorp.DbMap{Db: conn, Dialect: gorp.PostgresDialect{}}
	db.AddTableWithNameAndSchema(Audit{}, "workspace", "audit").SetKeys(true, "id")
	db.AddTableWithNameAndSchema(JobPG{}, "workspace", "jobs").SetKeys(true, "id")

	if debug {
		db.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	}

	return nil
}

func InitClickHouse() error {
	var (
		cfg = config.Config
		err error
	)

	dsn := fmt.Sprintf(
		"clickhouse://%s:%s@%s:%d/%s?dial_timeout=10s&read_timeout=20s",
		cfg.ClickHouse.User,
		cfg.ClickHouse.Password,
		cfg.ClickHouse.IP,
		cfg.ClickHouse.Port,
		cfg.ClickHouse.Name,
	)

	ch, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	//MigrateJobs()

	return nil
}
