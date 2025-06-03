package suresql

import (
	"fmt"

	orm "github.com/medatechnology/simpleorm"

	"github.com/medatechnology/goutil/filesystem"
	"github.com/medatechnology/goutil/print"
	"github.com/medatechnology/goutil/simplelog"
)

const (
	MIGRATION_DIRECTORY          = "migrations/"
	MIGRATION_UP_FILES_SIGNATURE = "_up.sql"
)

// This is more like migrating data from MIGRATION_DIRECTORY
// TODO: fix the printout to use metrics package so we can have the time elapsed information.
// Make sure to call this AFTER connect internal is called!! Because we need the DB connection already.
func InitDB(force bool) error {
	// If DB is already init, then do not run again
	if CurrentNode.Config.IsInitDone && !force {
		// simplelog.LogFormat("DB already initialized")
		return ErrDBInitializedAlready
	}

	simplelog.DEBUG_LEVEL = 1
	allUpFiles := filesystem.Dir(MIGRATION_DIRECTORY, MIGRATION_UP_FILES_SIGNATURE)
	fmt.Printf("\nMigration directory has %s files, proceed migration...",
		print.Colored(fmt.Sprintf("%d", len(allUpFiles)), print.ColorGreen))
	for _, ef := range allUpFiles {
		fContent := filesystem.More(MIGRATION_DIRECTORY + ef.Name())
		sqlCommands := orm.ConvertSQLCommands(fContent)
		fmt.Printf("Migrating file: %s - lines: %d - commands: %d",
			print.Colored(ef.Name(), print.ColorBlue), len(fContent), len(sqlCommands))
		// simplelog.LogFormat("Number of lines : %d", len(fContent))

		// DEBUG: print all the commands
		// for i, c := range sqlCommands {
		// 	fmt.Printf("%d:%s\n", i+1, c)
		// }
		// res, err := db.ExecManySQL(sqlCommands)

		res, err := CurrentNode.InternalConnection.ExecManySQL(sqlCommands)
		if err != nil {
			// NOTE: if one of the file has error, then cannot continue just return. Meaning could potentially initialized partially
			// TODO: create rollback functionality here.
			simplelog.LogErr(err, "cannot init migrate")
			return err
		}
		fmt.Printf("%d sql commands executed in : %sms\n", len(res), orm.SecondToMsString(orm.TotalTimeElapsedInSecond(res)))

		// simplelog.LogInfoAny("", 1, commands)
		// longSQL := strings.Join(fContent[:], "\n")
		// simplelog.LogThis(longSQL)

		// simplelog.LogThis(fContent...) // print in 1 line
		// for _, str := range fContent {
		// simplelog.LogFormat("%s", str)
		// }
	}
	res := CurrentNode.InternalConnection.ExecOneSQL("UPDATE " + CurrentNode.Config.TableName() + " SET is_init_done=true")
	if res.Error != nil {
		// NOTE: if one of the file has error, then cannot continue just return. Meaning could potentially initialized partially
		// TODO: create rollback functionality here.
		simplelog.LogErr(res.Error, "cannot update settings table")
		return res.Error
	}
	return nil
}
