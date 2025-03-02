package flags

import (
	"flag"
	"os"
)

type serverFlags struct {
	url         string
	accural     string
	databaseDSN string
}

func NewServerFlags() *serverFlags {
	return &serverFlags{}
}

var err error

func (f *serverFlags) SetServerFlags() {
	flag.StringVar(&f.url, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&f.accural, "r", "", "accural adress")
	flag.StringVar(&f.databaseDSN, "d", "postgres://gophermart:gophermart@localhost:5432/gophermart_db?sslmode=disable", "path to database")
}

func (f *serverFlags) GetServerURL() string {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.url = envRunAddr
	}

	return f.url
}

func (f *serverFlags) GetAccural() string {
	if envAccural := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccural != "" {
		f.accural = envAccural
	}
	return f.accural
}

func (f *serverFlags) GetDatabasePath() string {
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		f.databaseDSN = envDatabaseDSN
	}
	return f.databaseDSN
}
