package main

import (
	"github.com/medatechnology/suresql"
	"github.com/medatechnology/suresql/server"

	"github.com/medatechnology/goutil/simplelog"
)

// SureSQL BackEnd Service
func main() {
	err := suresql.ConnectInternal()
	if err != nil {
		// Cannot connect to DBMS, exit the app
		// TODO: maybe make this to just wait for 30sec then try again and keep looping like that
		// .     add the middleware CheckConnected and if false return "DBMS not connected" for all API
		simplelog.LogErrorStr("sureSQL", err, "Cannot connect to internal rqlite engine")
		return
	}

	// Prepare the SureSQL
	server := server.CreateServer(suresql.CurrentNode)

	suresql.CurrentNode.PrintWelcomePretty()
	// Start SureSQL server
	if err := server.Start(""); err != nil {
		simplelog.LogErrorStr("main", err, "cannot start SureSQL")
	}
}
