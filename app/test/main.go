package main

import (
	"fmt"
	"sync"

	"github.com/medatechnology/suresql"

	utils "github.com/medatechnology/goutil"
	"github.com/medatechnology/goutil/metrics"
	"github.com/medatechnology/goutil/simplelog"
)

func main() {
	simplelog.DEBUG_LEVEL = 1
	// config.LoadConfig("config/.env.dev", "config/.env.production")

	// err := suresql.ConnectInternal()
	// if err != nil {
	// 	simplelog.LogErrorStr("sureSQL", err, "Cannot connect to internal rqlite engine")
	// }

	utils.ReloadEnvEach(".env.dev")
	config := suresql.LoadDBMSConfigFromEnvironment()
	config.PrintDebug(true)

	db, err := suresql.NewDatabase(config)
	if err != nil {
		simplelog.LogErrorAny("Main", err, "Failed to connect to database")
	}
	simplelog.LogThis("Database loaded")
	fmt.Println(db.Status())
	// simplelog.LogAny(suresql.CurrentNode.InternalConnection.Status())

	// res := db.ExecOneSQL("DROP TABLE IF EXISTS users;")
	// if res.Err != nil {
	// 	simplelog.LogThis("cannot delete users table")
	// }
	// simplelog.LogThis("============= testing migration: BEGIN")
	// suresql.InitDB(false)
	// simplelog.LogThis("============= testing migration: DONE")

	// Test get Schema
	// schemas := suresql.CurrentNode.InternalConnection.GetSchema(true, false)
	schemas := db.GetSchema(true, false)
	simplelog.LogThis(fmt.Sprintf("Schemas contain %d tables and views", len(schemas)))
	for _, s := range schemas {
		s.PrintDebug(false)
	}
	// suresql.CurrentNode.PrintWelcomePretty()

	// TestBoxPrint()
	// simplelog.LogInfoAny("main", 1, schemas)

	// Test: Selectone - Works
	fmt.Println("Running select foo")
	res, err := db.SelectOne("foo")
	if err != nil {
		simplelog.LogErrorAny("Main", err, "Failed to connect to select")
	}
	fmt.Println("Done select foo")
	fmt.Println("Result = ", res)

	// Test: Many connection to rqlite
	// Run load test
	el := metrics.StartTimeIt("Running stress test concurrently...", 0)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db.SelectOneSQL("SELECT * FROM users LIMIT 1")
		}()
	}
	wg.Wait()
	metrics.StopTimeItPrint(el, "Done")

	// Test: Create Table and Insert Tble, works!
	// TestCreateTable(&suresql.CurrentNode.InternalConnection)
	// TestInsertTable(&suresql.CurrentNode.InternalConnection)

	// fmt.Println("Begin select one user")
	// TestSelectOne(db)
	// fmt.Println("Done select one user")

	// qr, err := db.conn.QueryOneParameterized(gorqlite.ParameterizedStatement{
	// 	Query: "SELECT * FROM users WHERE (age > ?) AND ((status = ?) OR (status = ?)) ORDER BY name ASC LIMIT 10",
	// 	Arguments: []interface{18, "active", "pending"},
	// })

	// fmt.Println("Begin select many user")
	// TestSelectMany(db)
	// fmt.Println("Done select many user")

}
