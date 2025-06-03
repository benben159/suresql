package server

import (
	"fmt"
	"time"

	"github.com/medatechnology/suresql"

	orm "github.com/medatechnology/simpleorm"
)

type AccessLogTable struct {
	ID            int       `json:"id,omitempty"              db:"id"`
	Username      string    `json:"username,omitempty"        db:"username"`
	ActionType    string    `json:"action_type,omitempty"     db:"action_type"`
	Occurred      time.Time `json:"occurred,omitempty"        db:"occurred"`
	Table         string    `json:"table_name,omitempty"      db:"table_name"`
	RawQuery      string    `json:"raw_query,omitempty"       db:"raw_query"`
	Result        string    `json:"result,omitempty"          db:"result"`
	ResultStatus  string    `json:"result_status,omitempty"   db:"result_status"`
	Error         string    `json:"error,omitempty"           db:"error"`
	Duration      float64   `json:"duration,omitempty"        db:"duration"`
	Method        string    `json:"method,omitempty"          db:"method"`
	NodeNumber    int       `json:"node_number,omitempty"     db:"node_number"`
	Note          string    `json:"note,omitempty"            db:"note"`
	Description   string    `json:"description,omitempty"     db:"description"`
	ClientIP      string    `json:"client_ip,omitempty"       db:"client_ip"`
	ClientBrowser string    `json:"client_browser,omitempty"  db:"client_browser"`
	ClientDevice  string    `json:"client_device,omitempty"   db:"client_device"`
}

func (l AccessLogTable) TableName() string {
	return "_access_logs"
}

// Save the logentry to log table
func (l *AccessLogTable) DBLogging(db *suresql.SureSQLDB) error {
	l.NodeNumber = suresql.CurrentNode.Config.NodeNumber
	l.Occurred = time.Now().UTC()
	// l.Username = db.Config.Username
	// return db.InsertOneTableStruct(l)
	return DBLogging(*db, *l)
}

func DBLogging(db suresql.SureSQLDB, entry AccessLogTable) error {

	// result := db.InsertOneTableStruct(entry, false)
	sql := `INSERT INTO %s (username, action_type, duration, result,
		result_status, error, method, node_number, note, description,
		client_ip, client_browser, client_device, 
		raw_query, table_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// fmt.Println("Occured:", entry.Occurred)
	// fmt.Println("Duration:", entry.Duration)
	// fmt.Println("Result:", entry.Result)
	params := orm.ParametereizedSQL{
		Query: fmt.Sprintf(sql, entry.TableName()),
		Values: []interface{}{
			entry.Username, entry.ActionType, entry.Duration, entry.Result,
			entry.ResultStatus, entry.Error, entry.Method, entry.NodeNumber, entry.Note, entry.Description,
			entry.ClientIP, entry.ClientBrowser, entry.ClientDevice,
			entry.RawQuery, entry.Table,
		},
	}

	result := db.ExecOneSQLParameterized(params)
	// if result.Err != nil {
	// 	fmt.Println("Error = ", result.Err)
	// } else {
	// 	fmt.Println("DB logged!")
	// }
	return result.Error
}
