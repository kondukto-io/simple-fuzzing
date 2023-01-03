package handlers

import "database/sql"

func MigrateDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE [users] (
ID TEXT NOT NULL PRIMARY KEY,
email TEXT NOT NULL,
name TEXT NOT NULL
);
   `)
	return err
}
