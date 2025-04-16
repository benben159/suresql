package suresql

import (
	"time"

	"github.com/medatechnology/simpleorm/rqlite"
)

var (
	ServerStartTime time.Time
)

// Making connection to internal DB
// This is where implementation selection happens, right now is only RQlite
func NewDatabase(conf SureSQLConfig) (SureSQLDB, error) {
	// TODO: FUTURE: maybe reading from environment then call the appropriate "connect" method to the internal
	// -             database. Currently is only RQLite, but maybe future can be postgres, mySQL etc.
	// conf.GenerateGoRQLiteURL()
	conf.GenerateRQLiteURL()

	config := rqlite.RqliteDirectConfig{
		URL:         conf.URL,
		Consistency: conf.Consistency,
		Username:    conf.Username,
		Password:    conf.Password,
		Timeout:     conf.HttpTimeout,
		RetryCount:  conf.MaxRetries,
	}
	SchemaTable = rqlite.SCHEMA_TABLE
	// TODO: make this read from environment
	CurrentNode.Status.DBMSDriver = "direct-rqlite"
	return rqlite.NewDatabase(config)
}
