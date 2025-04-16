package main

import (
	"fmt"

	orm "github.com/medatechnology/simpleorm"
	"github.com/medatechnology/suresql"

	"github.com/medatechnology/goutil/simplelog"
)

func TestSelectOne(db *suresql.SureSQLDB) {
	// Constructing a complex condition with grouping and ordering.
	ageLimit := 18
	condition := orm.Condition{
		Logic: "AND",
		Nested: []orm.Condition{
			{Field: "age", Operator: ">", Value: ageLimit},
			{
				Logic: "OR",
				Nested: []orm.Condition{
					{Field: "status", Operator: "=", Value: "active"},
					{Field: "status", Operator: "=", Value: "pending"},
				},
			},
		},
		// GroupBy: []string{"city"}, // Grouping by city
		OrderBy: []string{"name ASC"},
		Limit:   10,
		Offset:  0,
	}
	record, err := (*db).SelectOneWithCondition("users", &condition)
	if err != nil {
		simplelog.LogErrorAny("TestSelectOne", err, "Failed to select user")
		return
	}
	simplelog.LogFormat("User Data from DBRecord: %+v\n", record.Data)
}

func TestSelectMany(db *suresql.SureSQLDB) {
	// Constructing a complex condition with grouping and ordering.
	condition := orm.Condition{
		Logic: "AND",
		Nested: []orm.Condition{
			{Field: "age", Operator: ">", Value: 18},
			{
				Logic: "OR",
				Nested: []orm.Condition{
					{Field: "status", Operator: "=", Value: "active"},
					{Field: "status", Operator: "=", Value: "pending"},
				},
			},
		},
		// GroupBy: []string{"city"}, // Grouping by city
		OrderBy: []string{"name ASC"},
		Limit:   10,
		Offset:  0,
	}
	// For SelectMany usage:
	simplelog.LogThis("Begin query select many")
	records, err := (*db).SelectManyWithCondition("users", &condition)
	if err != nil {
		simplelog.LogErrorAny("TestSelectMany", err, "Failed to select user")
		return
	}
	fmt.Println("Successful query select many got rows:", len(records))
	fmt.Println("Records = ", len(records), " values:", records)
	// for idx, rec := range records {
	// 	simplelog.LogThis("Row ", fmt.Sprintf("%d", idx))
	// 	simplelog.LogFormat("User Data from DBRecord: %+v\n", rec.Data)
	// }
}
