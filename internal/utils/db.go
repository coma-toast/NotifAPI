package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DataModel struct {
	DB *sqlx.DB
}

type NotificationRow struct {
	PubID     string `db:"pub_id" json:"pub_id"`
	Date      string `db:"date" json:"date"`
	Source    string `db:"source" json:"source"`
	Interests string `db:"interests" json:"interests"`
	Title     string `db:"title" json:"title"`
	Message   string `db:"message" json:"message"`
	Metadata  string `db:"metadata" json:"metadata"`
}

func (d *DataModel) Init(location string) {
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		log.Fatal("Unable to create database directory", err)
	}

	d.DB = sqlx.MustConnect("sqlite3", location+"/data.db")

	schema := `CREATE TABLE IF NOT EXISTS notifications (
		pub_id text PRIMARY KEY,
		date TEXT DEFAULT CURRENT_TIMESTAMP,
		source TEXT,
		interests TEXT,
		title TEXT,
		message TEXT,
		metadata TEXT
	);`

	result := d.DB.MustExec(schema)
	fmt.Println("DB Initialized", result)
}

func (d *DataModel) AddNotification(pubID, source, title, message string, interests []string, metadata map[string]interface{}) (sql.Result, error) {
	insert := `INSERT INTO notifications 
	(
		pub_id,
		source,
		interests,
		title,
		message,
		metadata
	)
	VALUES 
	(
		?,
		?,
		?,
		?,
		?,
		?
	);`

	jsonInterests, err := json.Marshal(interests)
	if err != nil {
		return nil, err
	}

	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return d.DB.Exec(insert, pubID, source, jsonInterests, title, message, jsonMetadata)
}

func (d *DataModel) GetRecentNotifications(limit int) ([]NotificationRow, error) {
	notifications := make([]NotificationRow, 0)
	statement := `SELECT * FROM notifications ORDER BY date LIMIT ?`
	rows, err := d.DB.Queryx(statement, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row NotificationRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, row)
	}

	return notifications, nil

}
