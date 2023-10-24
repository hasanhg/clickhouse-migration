package config

type config struct {
	ClickHouse Database
	Database   Database
}

type Database struct {
	Driver          string
	IP              string
	Port            int
	User            string
	Password        string `json:"-"`
	Name            string
	ConnMaxIdle     int
	ConnMaxOpen     int
	ConnMaxLifetime int
	Debug           bool
}
