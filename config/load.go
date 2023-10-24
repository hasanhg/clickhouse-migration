package config

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/magiconair/properties"
)

// TODO: get rid of this global
var Config *config

func Load() (cfg *config, err error) {
	var path string
	for i, arg := range os.Args {
		if arg == "-v" {
			return nil, nil
		}
		path, err = parseCfg(os.Args, i)
		if err != nil {
			return nil, err
		}
		if path != "" {
			break
		}
	}
	p, err := loadProperties(path)
	if err != nil {
		return nil, err
	}

	Config, err = load(p)
	return Config, err
}

var errInvalidConfig = errors.New("invalid or missing path to config file")

func parseCfg(args []string, i int) (path string, err error) {
	if len(args) == 0 || i >= len(args) || !strings.HasPrefix(args[i], "-cfg") {
		return "", nil
	}
	arg := args[i]
	if arg == "-cfg" {
		if i >= len(args)-1 {
			return "", errInvalidConfig
		}
		return args[i+1], nil
	}

	if !strings.HasPrefix(arg, "-cfg=") {
		return "", errInvalidConfig
	}

	path = arg[len("-cfg="):]
	switch {
	case path == "":
		return "", errInvalidConfig
	case path[0] == '\'':
		path = strings.Trim(path, "'")
	case path[0] == '"':
		path = strings.Trim(path, "\"")
	}
	if path == "" {
		return "", errInvalidConfig
	}
	return path, nil
}

func loadProperties(path string) (p *properties.Properties, err error) {
	if path == "" {
		return properties.NewProperties(), nil
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return properties.LoadURL(path)
	}
	return properties.LoadFile(path, properties.UTF8)
}

func load(p *properties.Properties) (cfg *config, err error) {
	cfg = &config{}

	f := NewFlagSet(os.Args[0], flag.ExitOnError)

	// dummy values which were parsed earlier
	f.String("cfg", "", "Path or URL to config file")
	f.Bool("v", false, "Show version")

	// Database params
	f.StringVar(&cfg.Database.Driver, "postgres.driver", Default.Database.Driver, "Database driver to be used")
	f.StringVar(&cfg.Database.IP, "postgres.host", Default.Database.IP, "Database ip addr")
	f.IntVar(&cfg.Database.Port, "postgres.port", Default.Database.Port, "Database port")
	f.StringVar(&cfg.Database.User, "postgres.user", Default.Database.User, "Database username")
	f.StringVar(&cfg.Database.Password, "postgres.password", Default.Database.Password, "Database password")
	f.StringVar(&cfg.Database.Name, "postgres.db", Default.Database.Name, "Database name")
	f.BoolVar(&cfg.Database.Debug, "postgres.debug", Default.Database.Debug, "Database debug verbose mode")
	f.IntVar(&cfg.Database.ConnMaxIdle, "postgres.connmaxidle", Default.Database.ConnMaxIdle, "Maximum number of connections in the idle connection pool")
	f.IntVar(&cfg.Database.ConnMaxOpen, "postgres.connmaxopen", Default.Database.ConnMaxOpen, "Maximum number of open connections to the database")
	f.IntVar(&cfg.Database.ConnMaxLifetime, "postgres.connmaxlifetime", Default.Database.ConnMaxLifetime, "Maximum amount of time a connection may be reused")

	// ClickHouse params
	f.StringVar(&cfg.ClickHouse.IP, "clickhouse.ip", Default.ClickHouse.IP, "ClickHouse server ip addr")
	f.IntVar(&cfg.ClickHouse.Port, "clickhouse.port", Default.ClickHouse.Port, "ClickHouse server port")
	f.StringVar(&cfg.ClickHouse.User, "clickhouse.username", Default.ClickHouse.User, "Database username")
	f.StringVar(&cfg.ClickHouse.Password, "clickhouse.password", Default.ClickHouse.Password, "Database password")
	f.StringVar(&cfg.ClickHouse.Name, "clickhouse.name", Default.ClickHouse.Name, "Database name")

	// filter out -test flags
	var args []string
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-test.") {
			continue
		}
		args = append(args, a)
	}

	// parse configuration
	prefixes := []string{"", ""}
	if err := f.ParseFlags(args, os.Environ(), prefixes, p); err != nil {
		return nil, err
	}

	return cfg, nil
}
